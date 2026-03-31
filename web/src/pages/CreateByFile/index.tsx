import { useEffect, useMemo, useState } from 'react';
import {
  Alert,
  Button,
  Card,
  Collapse,
  Descriptions,
  Form,
  Input,
  Select,
  Space,
  Table,
  Tag,
  Typography,
  Upload,
  message,
} from 'antd';
import type { UploadProps } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { InboxOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { v4 as uuidv4 } from 'uuid';
import { getMcpConnectTokenByUUID, getMcpServerTools, uploadOpenapi } from '../../api';
import type {
  OpenapiUploadRequest,
  OpenapiUploadResponse,
  ToolInfo,
  ToolProtocolInfo,
} from '../../types';
import { AUTH_TYPE_LABELS, SERVER_STATUS_COLORS, SERVER_STATUS_LABELS } from '../../types';
import { buildMcpEndpoint } from '../../utils/mcp';

const { Dragger } = Upload;

function suffixFromFileName(name: string): string {
  const lower = name.toLowerCase();
  if (lower.endsWith('.json')) return 'json';
  if (lower.endsWith('.yaml') || lower.endsWith('.yml')) return 'yaml';
  return 'json';
}

function encodeBase64Utf8(value: string): string {
  const bytes = new TextEncoder().encode(value);
  let binary = '';
  bytes.forEach((byte) => {
    binary += String.fromCharCode(byte);
  });
  return btoa(binary);
}

const authSelectOptions = Object.keys(AUTH_TYPE_LABELS).map((value) => {
  const authValue = Number(value);
  return {
    value: authValue,
    label: AUTH_TYPE_LABELS[authValue],
  };
});

export default function CreateByFile() {
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const isAuth = Form.useWatch('isAuth', form);
  const authType = Number(isAuth);
  const [uploadFileName, setUploadFileName] = useState('');
  const [uploadSuffix, setUploadSuffix] = useState('json');
  const [submitting, setSubmitting] = useState(false);
  const [result, setResult] = useState<OpenapiUploadResponse | null>(null);
  const [connectToken, setConnectToken] = useState('');
  const [protocolTools, setProtocolTools] = useState<ToolProtocolInfo[]>([]);
  const [tokenLoading, setTokenLoading] = useState(false);
  const [protocolLoading, setProtocolLoading] = useState(false);

  const toolColumns: ColumnsType<ToolInfo> = [
    { title: '名称', dataIndex: 'name', key: 'name', ellipsis: true },
    { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
    { title: '方法', dataIndex: 'method', key: 'method', width: 100 },
    { title: '端点', dataIndex: 'endpoint', key: 'endpoint', ellipsis: true },
  ];

  const protocolItems = useMemo(
    () =>
      protocolTools.map((tool) => ({
        key: String(tool.id),
        label: tool.name || tool.fullName,
        children: (
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <Typography.Text type="secondary">
              {tool.description || '暂无描述'}
            </Typography.Text>
            <pre
              style={{
                margin: 0,
                padding: 12,
                background: '#f5f5f5',
                borderRadius: 6,
                overflow: 'auto',
                maxHeight: 320,
              }}
            >
              <code>{JSON.stringify(tool.inputSchema ?? {}, null, 2)}</code>
            </pre>
            <Alert
              type="warning"
              showIcon
              message="后端测试接口未实现"
              description="当前后端的 /v1/mcpServer/testMcpServerTool 没有返回成功响应，所以前端只能展示 MCP 协议工具列表，不能在这里执行真实测试。"
            />
          </Space>
        ),
      })),
    [protocolTools],
  );

  useEffect(() => {
    if (!result?.uuid) return;
    void fetchConnectToken(result.uuid);
    void fetchProtocolTools(result.uuid);
  }, [result?.uuid]);

  const beforeUpload: UploadProps['beforeUpload'] = (file) => {
    const name = file.name.toLowerCase();
    const ok = name.endsWith('.json') || name.endsWith('.yaml') || name.endsWith('.yml');
    if (!ok) {
      message.error('仅支持 .json、.yaml、.yml 文件');
      return Upload.LIST_IGNORE;
    }

    const reader = new FileReader();
    reader.onload = () => {
      const text = typeof reader.result === 'string' ? reader.result : '';
      form.setFieldValue('fileContent', text);
      setUploadFileName(file.name);
      setUploadSuffix(suffixFromFileName(file.name));
    };
    reader.onerror = () => {
      message.error('读取文件失败');
    };
    reader.readAsText(file, 'UTF-8');
    return false;
  };

  const onFinish = async (values: {
    name: string;
    description?: string;
    fileContent: string;
    isAuth: number;
    platformToken?: string;
    serviceToken?: string;
  }) => {
    setSubmitting(true);
    setResult(null);
    setConnectToken('');
    setProtocolTools([]);

    try {
      const payload: OpenapiUploadRequest = {
        uuid: uuidv4(),
        name: values.name,
        description: values.description ?? '',
        fileContent: encodeBase64Utf8(values.fileContent),
        suffix: uploadSuffix,
        isAuth: Number(values.isAuth) as OpenapiUploadRequest['isAuth'],
      };

      if (Number(values.isAuth) === 2 || Number(values.isAuth) === 4) {
        payload.platformToken = values.platformToken;
      }
      if (Number(values.isAuth) === 3 || Number(values.isAuth) === 4) {
        payload.serviceToken = values.serviceToken;
      }

      const res = await uploadOpenapi(payload);
      setResult(res.data.data);
      message.success('创建成功');
    } catch (e) {
      message.error(e instanceof Error ? e.message : '创建失败');
    } finally {
      setSubmitting(false);
    }
  };

  const fetchConnectToken = async (uuid: string) => {
    setTokenLoading(true);
    try {
      const res = await getMcpConnectTokenByUUID(uuid);
      const token = res.data.data?.connectToken ?? '';
      setConnectToken(token);
    } catch (e) {
      message.error(e instanceof Error ? e.message : '获取连接信息失败');
    } finally {
      setTokenLoading(false);
    }
  };

  const fetchProtocolTools = async (uuid: string) => {
    setProtocolLoading(true);
    try {
      const res = await getMcpServerTools(uuid);
      setProtocolTools(res.data.data ?? []);
    } catch (e) {
      message.error(e instanceof Error ? e.message : '获取 MCP 工具列表失败');
    } finally {
      setProtocolLoading(false);
    }
  };

  const mcpEndpoint = buildMcpEndpoint(connectToken);

  return (
    <div>
      <Typography.Title level={4} style={{ marginTop: 0 }}>
        通过 OpenAPI 文档创建 MCP Server
      </Typography.Title>

      <Form
        form={form}
        layout="vertical"
        style={{ maxWidth: 720 }}
        initialValues={{ isAuth: 1, description: '' }}
        onFinish={onFinish}
      >
        <Form.Item
          label="名称"
          name="name"
          rules={[{ required: true, message: '请输入名称' }]}
        >
          <Input placeholder="MCP Server 名称" />
        </Form.Item>

        <Form.Item label="描述" name="description">
          <Input.TextArea rows={3} placeholder="可选" />
        </Form.Item>

        <Form.Item
          label="OpenAPI 文件"
          required
          extra={uploadFileName ? `已选择：${uploadFileName}` : undefined}
        >
          <Dragger
            accept=".json,.yaml,.yml"
            maxCount={1}
            showUploadList={false}
            beforeUpload={beforeUpload}
          >
            <p className="ant-upload-drag-icon">
              <InboxOutlined />
            </p>
            <p className="ant-upload-text">点击或拖拽 OpenAPI 文件到此处</p>
            <p className="ant-upload-hint">支持 .json、.yaml、.yml</p>
          </Dragger>
        </Form.Item>

        <Form.Item
          name="fileContent"
          hidden
          rules={[{ required: true, message: '请上传 OpenAPI 文件' }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          label="鉴权类型"
          name="isAuth"
          rules={[{ required: true, message: '请选择鉴权类型' }]}
        >
          <Select options={authSelectOptions} placeholder="请选择" />
        </Form.Item>

        {(authType === 2 || authType === 4) && (
          <Form.Item
            label="平台 Token"
            name="platformToken"
            rules={[{ required: true, message: '请输入平台 Token' }]}
          >
            <Input placeholder="平台授权 Token" />
          </Form.Item>
        )}

        {(authType === 3 || authType === 4) && (
          <Form.Item
            label="服务 Token"
            name="serviceToken"
            rules={[{ required: true, message: '请输入服务 Token' }]}
          >
            <Input placeholder="Service 授权 Token" />
          </Form.Item>
        )}

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={submitting}>
            提交创建
          </Button>
        </Form.Item>
      </Form>

      {result ? (
        <Card title="创建结果" style={{ marginTop: 24 }}>
          <Space direction="vertical" size="large" style={{ width: '100%' }}>
            <Descriptions bordered size="small" column={1}>
              <Descriptions.Item label="ID">{result.id}</Descriptions.Item>
              <Descriptions.Item label="UUID">{result.uuid}</Descriptions.Item>
              <Descriptions.Item label="名称">{result.name}</Descriptions.Item>
              <Descriptions.Item label="描述">{result.description}</Descriptions.Item>
              <Descriptions.Item label="版本">{result.version}</Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color={SERVER_STATUS_COLORS[result.status]}>
                  {SERVER_STATUS_LABELS[result.status] ?? result.status}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="创建时间">{result.createdAt}</Descriptions.Item>
              <Descriptions.Item label="更新时间">{result.updatedAt}</Descriptions.Item>
            </Descriptions>

            <Card
              size="small"
              title="MCP 连接信息"
              extra={
                <Button type="link" onClick={() => navigate(`/server/${result.uuid}`)}>
                  查看详情页
                </Button>
              }
            >
              {tokenLoading ? (
                <Typography.Text type="secondary">正在获取 Connect Token...</Typography.Text>
              ) : connectToken ? (
                <Descriptions bordered size="small" column={1}>
                  <Descriptions.Item label="Connect Token">
                    <Typography.Text copyable>{connectToken}</Typography.Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="MCP Endpoint">
                    <Typography.Text copyable>{mcpEndpoint}</Typography.Text>
                  </Descriptions.Item>
                </Descriptions>
              ) : (
                <Typography.Text type="secondary">暂无连接信息</Typography.Text>
              )}
            </Card>

            <div>
              <Typography.Text strong>工具列表</Typography.Text>
              <Table<ToolInfo>
                style={{ marginTop: 12 }}
                rowKey={(row) => String(row.ID)}
                columns={toolColumns}
                dataSource={result.tools}
                pagination={false}
                scroll={{ x: true }}
                size="small"
              />
            </div>

            <Card size="small" title="MCP 协议工具列表">
              {protocolLoading ? (
                <Typography.Text type="secondary">正在加载 MCP 工具协议...</Typography.Text>
              ) : protocolTools.length ? (
                <Collapse items={protocolItems} />
              ) : (
                <Typography.Text type="secondary">当前没有可展示的 MCP 协议工具。</Typography.Text>
              )}
            </Card>
          </Space>
        </Card>
      ) : null}
    </div>
  );
}
