package biz

import (
	"context"
	"encoding/json"
	"errors"
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/conf"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/internal/mcp/config"
	"flow-bridge-mcp/internal/mcp/transformer"
	"flow-bridge-mcp/pkg/logger"
	"flow-bridge-mcp/pkg/tool"
	"fmt"
)

type OpenapiUseCase struct {
	log           *logger.Logger
	transformer   transformer.Transformer
	Tx            Transaction
	mcpServerRepo McpServerRepo
	mfUc          *McpFileUserCase
	mcpToolsRepo  McpToolsRepo
	conf          *conf.Conf
}

func NewOpenapiUserCase(
	transformer transformer.Transformer,
	log *logger.Logger,
	tx Transaction,
	mcpServerRepo McpServerRepo,
	mfUc *McpFileUserCase,
	mcpToolsRepo McpToolsRepo,
	conf *conf.Conf,
) *OpenapiUseCase {
	return &OpenapiUseCase{
		log:           log,
		transformer:   transformer,
		Tx:            tx,
		mcpServerRepo: mcpServerRepo,
		mfUc:          mfUc,
		mcpToolsRepo:  mcpToolsRepo,
		conf:          conf,
	}
}

func (o *OpenapiUseCase) Create(ctx context.Context, req *api.OpenapiUploadRequest) (resp *api.OpenapiUploadResponse, err error) {

	var decodeString []byte
	uuidStr := req.UUID
	if uuidStr == "" {
		o.log.ErrorWithContext(ctx, "UUID不能为空")
		return nil, fmt.Errorf("UUID不能为空")
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
		fileName := tool.FileNameByUUid() + "." + req.Suffix
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
	fmt.Printf("%+v", mcpInfo)
	//todo 如果文件存，则查询server服务
	serverUrls, err := json.Marshal(mcpInfo.Urls)
	if err != nil {
		o.log.ErrorWithContext(ctx, "mcpConfig.Urls json转换错误，err:%+v", err)
		return nil, err
	}
	var tools []string
	isHaveTools := false
	for _, val := range mcpInfo.Tools {
		if val.IsShow {
			isHaveTools = true
		}
		tools = append(tools, val.Name)
	}
	allTools, err := json.Marshal(tools)
	if err != nil {
		o.log.ErrorWithContext(ctx, "mcpConfig.Tools allowedTools json转换错误，err:%+v", err)
		return nil, err
	}
	haveTools := model.HaveToolsNo
	if isHaveTools {
		haveTools = model.HaveToolsYes
	}
	security, err := json.Marshal(mcpInfo.SecurityList)
	if err != nil {
		o.log.ErrorWithContext(ctx, "mcpConfig.Security json转换错误，err:%+v", err)
		return nil, err
	}

	err = o.Tx.ExecTx(ctx, func(ctx context.Context) error {
		var serverInfo = &model.McpServer{
			UUID:          req.UUID,
			Name:          req.Name,
			Description:   req.Description,
			Urls:          string(serverUrls),
			AllTools:      string(allTools),
			Version:       mcpInfo.Version,
			HaveTools:     haveTools,
			IsAuth:        model.IsAuthNo, //默认不开启权限控制，这里只是只平台的授权，接口的不需要开启
			ServiceToken:  "",
			PlatformToken: "",
			Security:      string(security),
			Status:        model.StatusHidden,
		}
		mcpServerId, err := o.mcpServerRepo.CreateWithTx(ctx, serverInfo)
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
			toolSecurity, err := json.Marshal(val.Security)
			if err != nil {
				o.log.ErrorWithContext(ctx, "mcpConfig.Tools.Security json转换错误，err:%+v", err)
				return fmt.Errorf("mcpConfig.Tools.Security json转换错误")
			}
			isAuth := model.IsAuthYes
			if val.SecurityLevel == config.SecurityLevelPublic {
				isAuth = model.IsAuthNo
			}
			inputSchema := ""
			if val.InputSchema != nil {
				inputSchemaByte, err := json.Marshal(val.InputSchema)
				if err != nil {
					o.log.ErrorWithContext(ctx, "mcpConfig.Tools.InputSchema json转换错误，err:%+v", err)
					return fmt.Errorf("mcpConfig.Tools.InputSchema json转换错误")
				}
				inputSchema = string(inputSchemaByte)
			}
			isShow := model.StatusDisplay
			if !val.IsShow {
				isShow = model.StatusHidden
			}

			var toolInfo = &model.McpTools{
				McpServerId:    mcpServerId,
				UUID:           req.UUID,
				Name:           val.Name,
				Description:    val.Description,
				McpServerType:  model.McpServerTypeOpenapi,
				Method:         val.Method,
				Endpoint:       val.Endpoint,
				Headers:        string(headers),
				Args:           string(args),
				RequestBody:    val.RequestBody,
				ResponseBody:   val.ResponseBody,
				InputSchema:    inputSchema,
				Annotations:    "", //暂时不做支持，如果需要可以考虑后期支持
				Security:       string(toolSecurity),
				IsAuth:         isAuth,
				AuthMode:       val.SecurityMode.String(),
				IsPlatformAuth: model.IsAuthNo, //默认不启用平台权限控制
				IsShow:         isShow,
			}
			mcpTools = append(mcpTools, toolInfo)

		}
		return o.mcpToolsRepo.CreateMcpToolsBatch(ctx, mcpServerId, req.UUID, tools, mcpTools)
	})
	if err != nil {
		o.log.ErrorWithContext(ctx, "创建McpServer失败，err:%+v", err)
		return nil, err
	}

	return
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
