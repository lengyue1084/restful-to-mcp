import { useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import {
  Alert,
  Button,
  Card,
  Descriptions,
  Form,
  Input,
  Select,
  Space,
  Typography,
  message,
} from 'antd';
import type {
  ArgConfig,
  AuthStatus,
  GetToolsInfoByUUIDResponse,
  TestMcpServerToolResponse,
} from '../../types';
import { AUTH_STATUS_LABELS, HTTP_METHODS } from '../../types';
import {
  createMcpServerTool,
  getToolsInfoByUUID,
  testMcpServerTool,
  updateMcpServerTool,
} from '../../api';

function parseArgs(raw: string): ArgConfig[] {
  const text = raw.trim();
  if (!text) return [];
  const parsed = JSON.parse(text) as unknown;
  if (!Array.isArray(parsed)) return [];

  return parsed.map((item: any) => ({
    name: String(item.name || ''),
    position: String(item.position || 'query'),
    required: !!item.required,
    type: String(item.type || 'string'),
    description: String(item.description || ''),
    default: String(item.default || ''),
    items: item.items,
    enum: Array.isArray(item.enum) ? item.enum.map((value: unknown) => String(value)) : undefined,
    explode: !!item.explode,
  }));
}

function parseArgumentsObject(raw: string): Record<string, unknown> {
  const text = raw.trim();
  if (!text) return {};
  const parsed = JSON.parse(text) as unknown;
  if (!parsed || Array.isArray(parsed) || typeof parsed !== 'object') {
    throw new Error('测试参数必须是合法的 JSON 对象');
  }
  return parsed as Record<string, unknown>;
}

function formatJson(value: unknown): string {
  if (value == null) return '{}';
  try {
    return JSON.stringify(value, null, 2);
  } catch {
    return String(value);
  }
}

export default function ToolEdit() {
  const { serverUUID, toolUUID } = useParams<{ serverUUID: string; toolUUID?: string }>();
  const [search] = useSearchParams();
  const navigate = useNavigate();
  const [form] = Form.useForm();

  const mode = search.get('mode') || 'edit';
  const detailMode = mode === 'detail';
  const isCreate = !toolUUID;

  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [testing, setTesting] = useState(false);
  const [toolInfo, setToolInfo] = useState<GetToolsInfoByUUIDResponse | null>(null);
  const [testArgsText, setTestArgsText] = useState('{}');
  const [testResult, setTestResult] = useState<TestMcpServerToolResponse | null>(null);

  const authOptions = useMemo(
    () =>
      (Object.keys(AUTH_STATUS_LABELS) as unknown as number[]).map((value) => ({
        value,
        label: AUTH_STATUS_LABELS[value],
      })),
    [],
  );

  useEffect(() => {
    if (!toolUUID) return;

    setLoading(true);
    void (async () => {
      try {
        const res = await getToolsInfoByUUID(toolUUID);
        const data = res.data.data;
        setToolInfo(data);
        form.setFieldsValue({
          name: data.name,
          description: data.description,
          method: data.method,
          path: data.path,
          isAuth: data.isAuth,
          isPlatformAuth: data.isPlatformAuth,
          argsText: JSON.stringify(data.args ?? [], null, 2),
        });
      } catch (error) {
        message.error(error instanceof Error ? error.message : '加载工具失败');
      } finally {
        setLoading(false);
      }
    })();
  }, [toolUUID, form]);

  const onFinish = async (values: any) => {
    if (!serverUUID) {
      message.error('缺少 serverUUID');
      return;
    }

    setSubmitting(true);
    try {
      const args = parseArgs(values.argsText || '[]');
      if (isCreate) {
        await createMcpServerTool({
          mcpServerUUID: serverUUID,
          name: values.name,
          description: values.description,
          method: values.method,
          path: values.path,
          platformAuth: values.isPlatformAuth as AuthStatus,
          toolAuth: values.isAuth as AuthStatus,
          args,
        });
        message.success('工具创建成功');
      } else {
        await updateMcpServerTool({
          uuid: toolUUID!,
          name: values.name,
          description: values.description,
          method: values.method,
          path: values.path,
          isAuth: values.isAuth as AuthStatus,
          isPlatformAuth: values.isPlatformAuth as AuthStatus,
          args,
        });
        message.success('工具更新成功');
      }
      navigate(`/server/${serverUUID}`);
    } catch (error) {
      message.error(error instanceof Error ? error.message : isCreate ? '工具创建失败' : '工具更新失败');
    } finally {
      setSubmitting(false);
    }
  };

  const handleTest = async () => {
    if (!toolUUID) return;

    setTesting(true);
    setTestResult(null);
    try {
      const argumentsObject = parseArgumentsObject(testArgsText);
      const res = await testMcpServerTool({
        uuid: toolUUID,
        arguments: argumentsObject,
      });
      setTestResult(res.data.data);
      message.success('工具测试成功');
    } catch (error) {
      message.error(error instanceof Error ? error.message : '工具测试失败');
    } finally {
      setTesting(false);
    }
  };

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <div>
        <Typography.Title level={4} style={{ marginTop: 0, marginBottom: 8 }}>
          {detailMode ? '工具详情' : isCreate ? '新增工具' : '编辑工具'}
        </Typography.Title>
        <Typography.Text type="secondary">
          当前后端支持 method、path、args、平台鉴权和接口鉴权的维护，并且已经补上真实工具测试接口。
        </Typography.Text>
      </div>

      <Card loading={loading}>
        <Space direction="vertical" size="middle" style={{ width: '100%' }}>
          <Alert
            type="info"
            showIcon
            message="参数录入说明"
            description="Path 中的 {id} 这类路径参数会由后端自动提取到工具 schema。args JSON 更适合补充 query、header、body 等额外入参。"
          />

          {toolInfo ? (
            <Descriptions bordered size="small" column={1}>
              <Descriptions.Item label="基础地址">{toolInfo.baseUrl || '暂无'}</Descriptions.Item>
              <Descriptions.Item label="Headers">
                <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>
                  <code>{formatJson(toolInfo.headers)}</code>
                </pre>
              </Descriptions.Item>
              <Descriptions.Item label="Security">
                <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>
                  <code>{formatJson(toolInfo.security)}</code>
                </pre>
              </Descriptions.Item>
            </Descriptions>
          ) : null}

          <Form
            form={form}
            layout="vertical"
            initialValues={{
              method: 'GET',
              isAuth: 1,
              isPlatformAuth: 1,
              argsText: '[]',
            }}
            onFinish={onFinish}
            disabled={detailMode}
          >
            <Form.Item
              name="name"
              label="名称"
              rules={[{ required: true, message: '请输入名称' }]}
            >
              <Input />
            </Form.Item>

            <Form.Item
              name="description"
              label="描述"
              rules={[{ required: true, message: '请输入描述' }]}
            >
              <Input.TextArea rows={3} />
            </Form.Item>

            <Form.Item
              name="method"
              label="HTTP Method"
              rules={[{ required: true, message: '请选择 Method' }]}
            >
              <Select options={HTTP_METHODS.map((method) => ({ label: method, value: method }))} />
            </Form.Item>

            <Form.Item
              name="path"
              label="Path"
              rules={[{ required: true, message: '请输入 Path' }]}
            >
              <Input placeholder="/users/{id}" />
            </Form.Item>

            <Form.Item
              name="isAuth"
              label="接口鉴权"
              rules={[{ required: true, message: '请选择接口鉴权' }]}
            >
              <Select options={authOptions} />
            </Form.Item>

            <Form.Item
              name="isPlatformAuth"
              label="平台鉴权"
              rules={[{ required: true, message: '请选择平台鉴权' }]}
            >
              <Select options={authOptions} />
            </Form.Item>

            <Form.Item
              name="argsText"
              label="请求入参配置 (JSON 数组)"
              extra='示例: [{"name":"page","position":"query","required":false,"type":"number","description":"页码","default":"1","explode":false}]'
              rules={[
                {
                  validator: async (_, value) => {
                    try {
                      parseArgs(value || '[]');
                    } catch {
                      throw new Error('请求入参必须是合法的 JSON 数组');
                    }
                  },
                },
              ]}
            >
              <Input.TextArea rows={12} />
            </Form.Item>

            <Space>
              {!detailMode ? (
                <Button type="primary" htmlType="submit" loading={submitting}>
                  {isCreate ? '创建' : '保存'}
                </Button>
              ) : null}
              <Button onClick={() => navigate(-1)}>返回</Button>
            </Space>
          </Form>
        </Space>
      </Card>

      {!isCreate ? (
        <Card
          title="工具测试"
          extra={(
            <Button type="primary" loading={testing} onClick={() => void handleTest()}>
              立即测试
            </Button>
          )}
        >
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <Alert
              type="warning"
              showIcon
              message="测试走真实下游请求"
              description="后端会读取当前 tool 的 path/args 配置，以及 server 上保存的 URL 和 serviceToken，直接向真实 HTTP 接口发请求。"
            />

            <Alert
              type="info"
              showIcon
              message="鉴权配置位置说明"
              description={(
                <Space direction="vertical" size={4} style={{ width: '100%' }}>
                  <Typography.Text>
                    1. 平台鉴权: 用于 MCP 网关访问，不是在本页测试参数里填写。外部客户端访问
                    <Typography.Text code style={{ marginInline: 4 }}>/gateway/:token/mcp</Typography.Text>
                    时，由网关校验
                    <Typography.Text code style={{ marginInline: 4 }}>x-mcp-platform-token</Typography.Text>
                    头，它的值来自 Server 配置中的
                    <Typography.Text code style={{ marginInline: 4 }}>platformToken</Typography.Text>。
                  </Typography.Text>
                  <Typography.Text>
                    2. 接口鉴权: 用于下游 REST 接口访问。请在 Server 配置页填写
                    <Typography.Text code style={{ marginInline: 4 }}>serviceToken</Typography.Text>
                    和安全配置
                    <Typography.Text code style={{ marginInline: 4 }}>Security Mode / Name / In / Scheme</Typography.Text>。
                    本页测试时会自动带上，不需要在下面的测试 JSON 里重复填写。
                  </Typography.Text>
                  <Typography.Text>
                    3. 普通业务参数: 比如 path/query/header/body 入参，才放在工具的
                    <Typography.Text code style={{ marginInline: 4 }}>args</Typography.Text>
                    配置和下方测试 JSON 里。固定请求头可以放在 Server 的
                    <Typography.Text code style={{ marginInline: 4 }}>Headers (JSON)</Typography.Text>
                    中。
                  </Typography.Text>
                </Space>
              )}
            />

            <div>
              <Typography.Text strong>测试参数 (JSON 对象)</Typography.Text>
              <Input.TextArea
                rows={10}
                value={testArgsText}
                onChange={(event) => setTestArgsText(event.target.value)}
                placeholder='{"id":"123","page":1}'
                style={{ marginTop: 8 }}
              />
            </div>

            {testResult ? (
              <Descriptions bordered size="small" column={1}>
                <Descriptions.Item label="工具名">{testResult.name}</Descriptions.Item>
                <Descriptions.Item label="请求方法">{testResult.requestMethod}</Descriptions.Item>
                <Descriptions.Item label="请求 URL">
                  <Typography.Text copyable>{testResult.requestUrl}</Typography.Text>
                </Descriptions.Item>
                <Descriptions.Item label="请求体">
                  <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>
                    <code>{formatJson(testResult.requestBody)}</code>
                  </pre>
                </Descriptions.Item>
                <Descriptions.Item label="响应 JSON">
                  <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>
                    <code>{formatJson(testResult.responseJson ?? {})}</code>
                  </pre>
                </Descriptions.Item>
                <Descriptions.Item label="响应文本">
                  <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>
                    <code>{testResult.responseText}</code>
                  </pre>
                </Descriptions.Item>
              </Descriptions>
            ) : null}
          </Space>
        </Card>
      ) : null}
    </Space>
  );
}
