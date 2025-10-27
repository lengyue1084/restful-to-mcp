package biz

import (
	"context"
	"encoding/json"
	"errors"
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/conf"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/internal/mcp/config"
	mcpServer "flow-bridge-mcp/internal/mcp/server"
	"flow-bridge-mcp/internal/mcp/transformer"
	"flow-bridge-mcp/internal/pkg/cache"
	"flow-bridge-mcp/pkg/const"
	"flow-bridge-mcp/pkg/logger"
	"flow-bridge-mcp/pkg/tool"
	"fmt"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

type OpenapiUseCase struct {
	log              *logger.Logger
	transformer      transformer.Transformer
	Tx               Transaction
	mcpServerRepo    McpServerRepo
	mfUc             *McpFileUserCase
	mcpToolsRepo     McpToolsRepo
	conf             *conf.Conf
	cache            *cache.MemoryCache
	mcpServerManager *mcpServer.McpServerManager
}

func NewOpenapiUserCase(
	transformer transformer.Transformer,
	log *logger.Logger,
	tx Transaction,
	mcpServerRepo McpServerRepo,
	mfUc *McpFileUserCase,
	mcpToolsRepo McpToolsRepo,
	conf *conf.Conf,
	cache *cache.MemoryCache,
	mcpServerManager *mcpServer.McpServerManager,
) *OpenapiUseCase {
	return &OpenapiUseCase{
		log:              log,
		transformer:      transformer,
		Tx:               tx,
		mcpServerRepo:    mcpServerRepo,
		mfUc:             mfUc,
		mcpToolsRepo:     mcpToolsRepo,
		conf:             conf,
		cache:            cache,
		mcpServerManager: mcpServerManager,
	}
}

func (o *OpenapiUseCase) Create(ctx context.Context, req *api.OpenapiUploadRequest) (resp *api.OpenapiUploadResponse, err error) {

	var decodeString []byte
	uuidStr := req.UUID
	if uuidStr == "" {
		return nil, fmt.Errorf("UUID不能为空")
	}
	//IsAuth        _const.AuthStatus `json:"isAuth" binding:"required"` //是否授权状态，这个状态是针对平台授权
	//ServiceToken  string            `json:"serviceToken" binding:"omitempty,min=1"`
	//PlatformToken string            `json:"platformToken" binding:"omitempty,min=1"`

	switch req.IsAuth {
	case _const.IsAuthServiceAuth:
		if req.ServiceToken == "" {
			return nil, fmt.Errorf("ServiceToken不能为空")
		}
	case _const.IsAuthPlatformAuth:
		if req.PlatformToken == "" {
			return nil, fmt.Errorf("PlatformToken不能为空")
		}
	case _const.IsAuthAllAuth:
		if req.ServiceToken == "" || req.PlatformToken == "" {
			return nil, fmt.Errorf("ServiceToken和PlatformToken不能为空")
		}
	}

	ctx = context.WithValue(ctx, "uuid", uuidStr)
	cleanedContent := tool.CleanBase64String(req.FileContent)
	if err := tool.ValidateBase64String(cleanedContent); err != nil {
		o.log.ErrorWithContext(ctx, "Base64字符串验证失败: %+v", err)
		return nil, err
	}
	decodeString, err = tool.TryMultipleBase64Decodings(cleanedContent)
	if err != nil {
		o.log.ErrorWithContext(ctx, "Base64解析失败，err:%+v", err)
		return nil, fmt.Errorf("base64 decode failed: %w", err)
	}
	contentMd5Str := tool.MD5(req.FileContent)

	//创建McpFile记录
	mcpFileInfo, err := o.mfUc.GetMcpFileInfoByMd5(ctx, contentMd5Str)
	if err != nil {
		o.log.ErrorWithContext(ctx, "获取McpFile失败，err:%+v", err)
		return nil, errors.New("获取McpFile失败")
	}
	if mcpFileInfo.ID == 0 {
		suffix := "yaml"
		if req.Suffix != "" {
			suffix = req.Suffix
		}
		fileName := tool.FileNameByUUid() + "." + suffix
		// 创建McpFile文件，写入记录
		err := o.mfUc.CreateMcpFile(ctx, fileName, req.Name, req.Suffix, contentMd5Str, req.Description)
		if err != nil {
			o.log.ErrorWithContext(ctx, "创建McpFile失败，err:%+v", err)
			return nil, errors.New("创建McpFile失败")
		}

		path := o.conf.Conf.GetString("file.path")
		//写入文件中,md5表明是否变化
		_, err = tool.WriteFile(path, fileName, decodeString)
		if err != nil {
			o.log.ErrorWithContext(ctx, "写入文件失败，err:%+v", err)
			return nil, err
		}
	} else {
		err = o.mfUc.UpdateMcpFile(ctx, int64(mcpFileInfo.ID), req.Name, req.Suffix, contentMd5Str, req.Description)
		if err != nil {
			o.log.ErrorWithContext(ctx, "更新McpFile失败，err:%+v", err)
			return nil, errors.New("更新McpFile失败")
		}
	}

	mcpInfo, err := o.transformer.Convert(ctx, decodeString)
	if err != nil {
		o.log.ErrorWithContext(ctx, "数据转换错误，err:%+v", err)
		return nil, err
	}
	serverUrls, err := json.Marshal(mcpInfo.Urls)
	if err != nil {
		o.log.ErrorWithContext(ctx, "mcpConfig.Urls json转换错误，err:%+v", err)
		return nil, err
	}
	var tools []string
	haveTools := _const.HaveToolsNo
	for _, val := range mcpInfo.Tools {
		if val.IsShow {
			haveTools = _const.HaveToolsYes
		}
		tools = append(tools, val.Name)
	}
	allTools, err := json.Marshal(tools)
	if err != nil {
		o.log.ErrorWithContext(ctx, "mcpConfig.Tools allowedTools json转换错误，err:%+v", err)
		return nil, err
	}
	security := ""
	if mcpInfo.SecurityList != nil {
		securityByte, err := json.Marshal(mcpInfo.SecurityList)
		if err != nil {
			o.log.ErrorWithContext(ctx, "mcpConfig.Security json转换错误，err:%+v", err)
			return nil, err
		}
		security = string(securityByte)
	}

	var (
		serialNumber = ""
		maxRetries   = _const.CommonRetryTimes
		retryCount   = 0
	)
	for retryCount < maxRetries {
		serialNumber = tool.RandStringWithLowercaseAndDigits(6)
		number, err := o.mcpServerRepo.GetCountMcpServerInfoBySerialNumber(ctx, serialNumber)
		if err != nil {
			o.log.ErrorWithContext(ctx, "查询序列号失败，err:%+v", err)
			return nil, err
		}
		if number == 0 {
			// 找到唯一的序列号
			break
		}
		retryCount++
		o.log.WarnWithContext(ctx, "序列号已存在，重新生成，尝试次数: %d", retryCount)
	}

	// 检查是否成功生成唯一序列号
	if retryCount >= maxRetries || serialNumber == "" {
		o.log.ErrorWithContext(ctx, "生成唯一序列号失败，已达到最大重试次数: %d", maxRetries)
		return nil, fmt.Errorf("生成唯一序列号失败，请稍后重试")
	}

	var mcpServerId = int64(0)
	err = o.Tx.ExecTx(ctx, func(ctx context.Context) (err error) {
		var serverInfo = &model.McpServer{
			UUID:          req.UUID,
			Name:          req.Name,
			Description:   req.Description,
			Urls:          string(serverUrls),
			AllTools:      string(allTools),
			Version:       mcpInfo.Version,
			HaveTools:     haveTools,
			IsAuth:        req.IsAuth, //默认不开启权限控制，这里只是只平台的授权，接口的不需要开启
			ServiceToken:  req.ServiceToken,
			PlatformToken: req.PlatformToken,
			Security:      security,
			Status:        _const.ServerNotSetToken,
			SerialNumber:  serialNumber,
			McpServerType: _const.McpServerTypeOpenapi,
			Source:        _const.SourceTypeFile,
		}
		mcpServerId, err = o.mcpServerRepo.CreateWithTx(ctx, serverInfo)
		if err != nil {
			o.log.ErrorWithContext(ctx, "CreateWithTx,创建McpServer失败，err:%+v", err)
			return fmt.Errorf("创建McpServer失败")
		}
		var mcpTools []*model.McpTools
		for _, val := range mcpInfo.Tools {
			headers, err := json.Marshal(val.Headers)
			if err != nil {
				o.log.ErrorWithContext(ctx, "mcpConfig.Tools.Headers json转换错误，err:%+v", err)
				return fmt.Errorf("mcpConfig.Tools.Headers json转换错误")
			}
			args, err := json.Marshal(val.Args)
			if err != nil {
				o.log.ErrorWithContext(ctx, "mcpConfig.Tools.Args json转换错误，err:%+v", err)
				return fmt.Errorf("mcpConfig.Tools.Args json转换错误")
			}
			var toolSecurity = ""
			if val.Security != nil {
				toolSecurityByte, err := json.Marshal(val.Security)
				if err != nil {
					o.log.ErrorWithContext(ctx, "mcpConfig.Tools.Security json转换错误，err:%+v", err)
					return fmt.Errorf("mcpConfig.Tools.Security json转换错误")
				}
				toolSecurity = string(toolSecurityByte)
			}

			isAuth := _const.IsAuthYes
			if val.SecurityLevel == config.SecurityLevelPublic {
				isAuth = _const.IsAuthNo
			}
			toolSchema := ""
			if val.ToolSchema != nil {
				inputSchemaByte, err := json.Marshal(val.ToolSchema)
				if err != nil {
					o.log.ErrorWithContext(ctx, "mcpConfig.Tools.InputSchema json转换错误，err:%+v", err)
					return fmt.Errorf("mcpConfig.Tools.InputSchema json转换错误")
				}
				toolSchema = string(inputSchemaByte)
			}
			annotations := ""
			if val.Annotations != nil {
				annotationsByte, err := json.Marshal(val.Annotations)
				if err != nil {
					o.log.ErrorWithContext(ctx, "mcpConfig.Tools.Annotations json转换错误，err:%+v", err)
					return fmt.Errorf("mcpConfig.Tools.Annotations json转换错误")
				}
				annotations = string(annotationsByte)
			}
			isShow := _const.StatusDisplay
			if !val.IsShow {
				isShow = _const.StatusHidden
			}

			var toolInfo = &model.McpTools{
				UUID:           tool.NewUUID(),
				McpServerId:    mcpServerId,
				McpServerUUID:  req.UUID,
				Name:           val.Name,
				Description:    val.Description,
				McpServerType:  _const.McpServerTypeOpenapi,
				Method:         val.Method,
				Endpoint:       val.Endpoint,
				Headers:        string(headers),
				Args:           string(args),
				RequestBody:    val.RequestBody,
				ResponseBody:   val.ResponseBody,
				ToolSchema:     toolSchema,
				Annotations:    annotations, //暂时不做支持，如果需要可以考虑后期支持
				Security:       toolSecurity,
				IsAuth:         isAuth,
				AuthMode:       val.SecurityMode.String(),
				IsPlatformAuth: _const.IsAuthNo, //默认不启用平台权限控制
				IsShow:         isShow,
				SerialNumber:   serialNumber,
				IsRepeat:       _const.CommonStatusNo,
			}
			mcpTools = append(mcpTools, toolInfo)

		}
		err = o.mcpToolsRepo.CreateMcpToolsBatch(ctx, mcpServerId, req.UUID, tools, mcpTools)
		if err != nil {
			o.log.ErrorWithContext(ctx, "创建McpTools失败，err:%+v", err)
			return err
		}
		return
	})
	if err != nil {
		o.log.ErrorWithContext(ctx, "创建McpServer失败，err:%+v", err)
		return nil, err
	}
	mcpServerInfo, err := o.mcpServerRepo.GetMcpServerInfoByID(ctx, mcpServerId)
	if err != nil {
		o.log.ErrorWithContext(ctx, "获取McpServer失败，err:%+v", err)
		return nil, err
	}
	var urls []string
	err = json.Unmarshal([]byte(mcpServerInfo.Urls), &urls)
	if err != nil {
		o.log.ErrorWithContext(ctx, "mcpServerInfo.Urls json转换错误，err:%+v", err)
		return nil, err
	}
	var toolsList []*api.ToolInfo
	for _, val := range mcpServerInfo.Tools {
		var toolSchema *protocol.InputSchema
		if val.ToolSchema != "" {
			err = json.Unmarshal([]byte(val.ToolSchema), &toolSchema)
			if err != nil {
				o.log.ErrorWithContext(ctx, "mcpServerInfo.ToolSchema json转换错误，err:%+v", err)
				return nil, err
			}
		}
		toolsList = append(toolsList, &api.ToolInfo{
			ID:             val.ID,
			McpServerId:    val.McpServerId,
			UUID:           val.McpServerUUID,
			Name:           val.Name,
			Description:    val.Description,
			Method:         val.Method,
			Endpoint:       val.Endpoint,
			Headers:        val.Headers,
			Args:           val.Args,
			RequestBody:    val.RequestBody,
			ResponseBody:   val.ResponseBody,
			ToolSchema:     toolSchema,
			Annotations:    val.Annotations,
			Security:       val.Security,
			IsAuth:         val.IsAuth,
			IsShow:         val.IsShow,
			IsPlatformAuth: val.IsPlatformAuth,
			CreatedAt:      val.CreatedAt,
			UpdatedAt:      val.UpdatedAt,
		})
	}
	resp = &api.OpenapiUploadResponse{
		ID:          mcpServerInfo.ID,
		UUID:        mcpServerInfo.UUID,
		Name:        mcpServerInfo.Name,
		Description: mcpServerInfo.Description,
		Urls:        urls,
		AllTools:    tools,
		Version:     mcpServerInfo.Version,
		Tools:       toolsList,
		CreatedAt:   mcpServerInfo.CreatedAt,
		UpdatedAt:   mcpServerInfo.UpdatedAt,
		Status:      mcpServerInfo.Status,
	}
	return resp, nil
}

func (o *OpenapiUseCase) UpdateForAuth(ctx context.Context, req *api.OpenapiUpdateForAuthRequest) (resp *api.OpenapiUpdateForAuthResponse, err error) {

	//更新老的接口到缓存
	o.UpdateToolsForOldCache(ctx)
	err = o.Tx.ExecTx(ctx, func(ctx context.Context) error {
		err = o.mcpServerRepo.UpdateMcpServerForAuthWithTx(ctx, req.UUID, req.IsAuth, req.ServiceToken, req.PlatformToken)
		if err != nil {
			o.log.ErrorWithContext(ctx, "更新server token 失败,err:%+v", err)
			return fmt.Errorf("更新失败，err:%+v", err)
		}
		if len(req.Tools) == 0 {
			return nil
		}
		err = o.mcpToolsRepo.UpdateToolsForAuthWithTx(ctx, req.UUID, req.Tools)
		if err != nil {
			o.log.ErrorWithContext(ctx, "更新tool token 失败,err:%+v", err)
			return fmt.Errorf("更新失败，err:%+v", err)
		}
		return nil
	})
	if err != nil {
		o.log.ErrorWithContext(ctx, "更新token 信息失败,err:%+v", err)
		return nil, fmt.Errorf("更新失败，err:%+v", err)
	}
	go func(ctx2 context.Context) {
		defer func() {
			if err := recover(); err != nil {
				o.log.ErrorWithContext(ctx2, "panic: %+v", err)
			}
		}()
		o.UpdateToolsForCache(ctx2)
		o.mcpServerManager.RegisterToolFromCache()
	}(ctx)

	return
}
func (o *OpenapiUseCase) UpdateToolsForOldCache(ctx context.Context) {
	mcpServerInfo, err := o.mcpServerRepo.GetMcpServerInfoWithAllTools(ctx)
	if err != nil {
		o.log.ErrorWithContext(context.Background(), "UpdateToolsForOldCache GetMcpServerInfoWithTools error: %v", err)
		return
	}
	o.cache.StoreMcpServer(cache.OldMcpValue, mcpServerInfo)
}
func (o *OpenapiUseCase) UpdateToolsForCache(ctx context.Context) {

	mcpServerInfo, err := o.mcpServerRepo.GetMcpServerInfoWithAllTools(ctx)
	if err != nil {
		o.log.ErrorWithContext(context.Background(), "GetMcpServerInfoWithTools error: %v", err)
		return
	}
	o.cache.StoreMcpServer(cache.NewMcpValue, mcpServerInfo)
}

func (o *OpenapiUseCase) ValidateBase64String(ctx context.Context, fileContent string) bool {
	// 清理Base64字符串
	cleanedContent := tool.CleanBase64String(fileContent)
	o.log.ErrorWithContext(ctx, "UUID不能为空")
	// 验证Base64字符串
	if err := tool.ValidateBase64String(cleanedContent); err != nil {
		o.log.ErrorWithContext(ctx, "Base64字符串验证失败: %+v", err)
		return false
	}
	return true

}
