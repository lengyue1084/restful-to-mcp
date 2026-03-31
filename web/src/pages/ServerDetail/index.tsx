import { useCallback, useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Alert,
  Button,
  Card,
  Collapse,
  Descriptions,
  Form,
  Input,
  Modal,
  Popconfirm,
  Space,
  Spin,
  Switch,
  Table,
  Tabs,
  Tag,
  Typography,
  message,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
  deleteMcpServerByUUID,
  getMcpConnectTokenByUUID,
  getMcpServerInfoByUUID,
  getMcpServerTools,
  getMcpServerToolsByUUID,
  updateForAuth,
  updateMcpServerByUUID,
} from '../../api';
import type {
  CommonToolItemInfo,
  GetMcpServerInfoByUUIDResponse,
  ToolProtocolInfo,
} from '../../types';
import {
  AUTH_STATUS_LABELS,
  AUTH_TYPE_LABELS,
  COMMON_STATUS_LABELS,
  SERVER_STATUS_COLORS,
  SERVER_STATUS_LABELS,
  SOURCE_TYPE_LABELS,
  STATUS_LABELS,
} from '../../types';
import { buildMcpEndpoint } from '../../utils/mcp';

function methodTagColor(method: string): string {
  const value = method.toUpperCase();
  if (value === 'GET') return 'green';
  if (value === 'POST') return 'blue';
  if (value === 'PUT') return 'orange';
  if (value === 'DELETE') return 'red';
  return 'default';
}

function formatJson(value: unknown): string {
  if (value == null) return '{}';
  try {
    return JSON.stringify(value, null, 2);
  } catch {
    return String(value);
  }
}

function buildCurlCommand(endpoint: string, payload: Record<string, unknown>): string {
  if (!endpoint) return '';
  const body = JSON.stringify(payload).replace(/"/g, '\\"');
  return `curl.exe -X POST "${endpoint}" -H "Content-Type: application/json" -H "Accept: application/json, text/event-stream" -d "${body}"`;
}

function CodeBlock({ value }: { value: string }) {
  return (
    <pre
      style={{
        margin: 0,
        padding: 12,
        background: '#f5f5f5',
        borderRadius: 6,
        overflow: 'auto',
        maxHeight: 360,
        whiteSpace: 'pre-wrap',
        wordBreak: 'break-all',
      }}
    >
      <code>{value}</code>
    </pre>
  );
}

export default function ServerDetail() {
  const { uuid } = useParams<{ uuid: string }>();
  const navigate = useNavigate();
  const [form] = Form.useForm<{ name: string; description: string }>();

  const [serverInfo, setServerInfo] = useState<GetMcpServerInfoByUUIDResponse | null>(null);
  const [tools, setTools] = useState<CommonToolItemInfo[]>([]);
  const [protocolTools, setProtocolTools] = useState<ToolProtocolInfo[]>([]);
  const [connectToken, setConnectToken] = useState('');

  const [pageLoading, setPageLoading] = useState(true);
  const [toolsLoading, setToolsLoading] = useState(false);
  const [protocolLoading, setProtocolLoading] = useState(false);
  const [editSubmitting, setEditSubmitting] = useState(false);
  const [tokenLoading, setTokenLoading] = useState(false);
  const [deleteLoading, setDeleteLoading] = useState(false);
  const [batchAuthSubmitting, setBatchAuthSubmitting] = useState(false);

  const [editOpen, setEditOpen] = useState(false);
  const [batchAuthOpen, setBatchAuthOpen] = useState(false);
  const [authDraft, setAuthDraft] = useState<Record<number, boolean>>({});

  const loadServer = useCallback(async () => {
    if (!uuid) return;
    setPageLoading(true);
    try {
      const res = await getMcpServerInfoByUUID(uuid);
      setServerInfo(res.data.data);
    } catch (error) {
      message.error(error instanceof Error ? error.message : '获取服务信息失败');
    } finally {
      setPageLoading(false);
    }
  }, [uuid]);

  const loadTools = useCallback(async () => {
    if (!uuid) return;
    setToolsLoading(true);
    try {
      const res = await getMcpServerToolsByUUID(uuid);
      setTools(res.data.data?.tools ?? []);
    } catch (error) {
      message.error(error instanceof Error ? error.message : '获取工具列表失败');
    } finally {
      setToolsLoading(false);
    }
  }, [uuid]);

  const loadProtocolTools = useCallback(async () => {
    if (!uuid) return;
    setProtocolLoading(true);
    try {
      const res = await getMcpServerTools(uuid);
      setProtocolTools(res.data.data ?? []);
    } catch (error) {
      message.error(error instanceof Error ? error.message : '获取 MCP 协议工具列表失败');
    } finally {
      setProtocolLoading(false);
    }
  }, [uuid]);

  useEffect(() => {
    if (!uuid) {
      message.error('缺少服务 UUID');
      return;
    }

    void loadServer();
    void loadTools();
    void loadProtocolTools();
  }, [uuid, loadServer, loadTools, loadProtocolTools]);

  const openEdit = () => {
    if (!serverInfo) return;
    form.setFieldsValue({
      name: serverInfo.name,
      description: serverInfo.description,
    });
    setEditOpen(true);
  };

  const submitEdit = async () => {
    if (!uuid) return;
    const values = await form.validateFields();
    setEditSubmitting(true);
    try {
      await updateMcpServerByUUID({
        uuid,
        name: values.name,
        description: values.description,
      });
      message.success('保存成功');
      setEditOpen(false);
      await loadServer();
    } catch (error) {
      message.error(error instanceof Error ? error.message : '保存失败');
    } finally {
      setEditSubmitting(false);
    }
  };

  const fetchConnectToken = async () => {
    if (!uuid) return;
    setTokenLoading(true);
    try {
      const res = await getMcpConnectTokenByUUID(uuid);
      setConnectToken(res.data.data?.connectToken ?? '');
      message.success('已生成 connectToken 连接信息');
    } catch (error) {
      message.error(error instanceof Error ? error.message : '获取连接信息失败');
    } finally {
      setTokenLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!uuid) return;
    setDeleteLoading(true);
    try {
      await deleteMcpServerByUUID(uuid);
      message.success('已删除');
      navigate('/', { replace: true });
    } catch (error) {
      message.error(error instanceof Error ? error.message : '删除失败');
    } finally {
      setDeleteLoading(false);
    }
  };

  const openBatchAuth = () => {
    const nextDraft: Record<number, boolean> = {};
    tools.forEach((tool) => {
      nextDraft[tool.id] = tool.isAuth === 2;
    });
    setAuthDraft(nextDraft);
    setBatchAuthOpen(true);
  };

  const submitBatchAuth = async () => {
    if (!uuid) return;
    setBatchAuthSubmitting(true);
    try {
      await updateForAuth({
        uuid,
        tools: tools.map((tool) => ({
          id: tool.id,
          isAuth: (authDraft[tool.id] ? 2 : 1) as 1 | 2,
        })),
      });
      message.success('接口鉴权已更新');
      setBatchAuthOpen(false);
      await loadServer();
      await loadTools();
      await loadProtocolTools();
    } catch (error) {
      message.error(error instanceof Error ? error.message : '接口鉴权更新失败');
    } finally {
      setBatchAuthSubmitting(false);
    }
  };

  const uuidEndpoint = useMemo(() => buildMcpEndpoint(uuid ?? ''), [uuid]);
  const connectTokenEndpoint = useMemo(() => buildMcpEndpoint(connectToken), [connectToken]);

  const firstToolName = useMemo(() => {
    const visibleProtocolTool = protocolTools.find((item) => item.isShow === 2);
    if (visibleProtocolTool?.name) return visibleProtocolTool.name;
    if (visibleProtocolTool?.fullName) return visibleProtocolTool.fullName;

    const visibleTool = tools.find((item) => item.isShow === 2);
    if (visibleTool?.name) return visibleTool.name;

    return '<tool_name>';
  }, [protocolTools, tools]);

  const toolsListPayload = useMemo(
    () => ({ jsonrpc: '2.0', id: 1, method: 'tools/list', params: {} }),
    [],
  );

  const buildToolsCallPayload = useCallback(
    (toolName: string) => ({
      jsonrpc: '2.0',
      id: 2,
      method: 'tools/call',
      params: {
        name: toolName,
        arguments: {},
      },
    }),
    [],
  );

  const uuidToolsListCurl = useMemo(
    () => buildCurlCommand(uuidEndpoint, toolsListPayload),
    [uuidEndpoint, toolsListPayload],
  );
  const uuidToolsCallCurl = useMemo(
    () => buildCurlCommand(uuidEndpoint, buildToolsCallPayload(firstToolName)),
    [buildToolsCallPayload, firstToolName, uuidEndpoint],
  );
  const connectTokenToolsListCurl = useMemo(
    () => buildCurlCommand(connectTokenEndpoint, toolsListPayload),
    [connectTokenEndpoint, toolsListPayload],
  );
  const connectTokenToolsCallCurl = useMemo(
    () => buildCurlCommand(connectTokenEndpoint, buildToolsCallPayload(firstToolName)),
    [buildToolsCallPayload, connectTokenEndpoint, firstToolName],
  );

  const toolColumns: ColumnsType<CommonToolItemInfo> = useMemo(
    () => [
      { title: '名称', dataIndex: 'name', key: 'name', width: 180, ellipsis: true },
      { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
      {
        title: '方法',
        dataIndex: 'method',
        key: 'method',
        width: 100,
        render: (value: string) => <Tag color={methodTagColor(value)}>{value}</Tag>,
      },
      { title: '端点', dataIndex: 'endpoint', key: 'endpoint', ellipsis: true },
      {
        title: '接口鉴权',
        dataIndex: 'isAuth',
        key: 'isAuth',
        width: 100,
        render: (value: number) => AUTH_STATUS_LABELS[value] ?? value,
      },
      {
        title: '平台鉴权',
        dataIndex: 'isPlatformAuth',
        key: 'isPlatformAuth',
        width: 100,
        render: (value: number) => AUTH_STATUS_LABELS[value] ?? value,
      },
      {
        title: '显示',
        dataIndex: 'isShow',
        key: 'isShow',
        width: 80,
        render: (value: number) => STATUS_LABELS[value] ?? value,
      },
      {
        title: '重名',
        dataIndex: 'isRepeat',
        key: 'isRepeat',
        width: 80,
        render: (value: number) => COMMON_STATUS_LABELS[value] ?? value,
      },
      {
        title: '操作',
        key: 'actions',
        width: 180,
        fixed: 'right',
        render: (_, record) => (
          <Space size="small">
            <Button type="link" size="small" onClick={() => navigate(`/server/${uuid}/tool/${record.uuid}?mode=edit`)}>
              编辑
            </Button>
            <Button type="link" size="small" onClick={() => navigate(`/server/${uuid}/tool/${record.uuid}?mode=detail`)}>
              详情
            </Button>
          </Space>
        ),
      },
    ],
    [navigate, uuid],
  );

  const batchAuthColumns: ColumnsType<CommonToolItemInfo> = useMemo(
    () => [
      { title: '名称', dataIndex: 'name', key: 'name', ellipsis: true },
      {
        title: '启用接口鉴权',
        key: 'isAuth',
        width: 180,
        render: (_, record) => (
          <Switch
            checked={!!authDraft[record.id]}
            onChange={(checked) => setAuthDraft((state) => ({ ...state, [record.id]: checked }))}
          />
        ),
      },
    ],
    [authDraft],
  );

  const protocolItems = useMemo(
    () =>
      protocolTools.map((tool) => ({
        key: String(tool.id),
        label: tool.name || tool.fullName,
        children: (
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <Typography.Text type="secondary">{tool.description || '暂无描述'}</Typography.Text>
            <CodeBlock value={JSON.stringify(tool.inputSchema ?? {}, null, 2)} />
          </Space>
        ),
      })),
    [protocolTools],
  );

  const tabItems = [
    {
      key: 'basic',
      label: '基本信息',
      children: (
        <Spin spinning={pageLoading}>
          {serverInfo ? (
            <Space direction="vertical" size="large" style={{ width: '100%' }}>
              <Space wrap>
                <Button type="primary" onClick={openEdit}>
                  编辑名称和描述
                </Button>
                {serverInfo.source === 2 ? (
                  <Button onClick={() => navigate(`/create/form/${serverInfo.uuid}`)}>编辑完整配置</Button>
                ) : null}
                <Button loading={tokenLoading} onClick={() => void fetchConnectToken()}>
                  生成 connectToken 连接信息
                </Button>
                <Popconfirm
                  title="确定删除这个 MCP Server 吗？"
                  okText="删除"
                  cancelText="取消"
                  okButtonProps={{ loading: deleteLoading, danger: true }}
                  onConfirm={() => void handleDelete()}
                >
                  <Button danger loading={deleteLoading}>
                    删除
                  </Button>
                </Popconfirm>
              </Space>

              <Descriptions bordered size="small" column={1}>
                <Descriptions.Item label="UUID">{serverInfo.uuid}</Descriptions.Item>
                <Descriptions.Item label="名称">{serverInfo.name}</Descriptions.Item>
                <Descriptions.Item label="描述">{serverInfo.description || '暂无'}</Descriptions.Item>
                <Descriptions.Item label="版本">{serverInfo.version || '暂无'}</Descriptions.Item>
                <Descriptions.Item label="状态">
                  <Tag color={SERVER_STATUS_COLORS[serverInfo.status]}>
                    {SERVER_STATUS_LABELS[serverInfo.status] ?? serverInfo.status}
                  </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="认证类型">
                  {AUTH_TYPE_LABELS[serverInfo.isAuth] ?? serverInfo.isAuth}
                </Descriptions.Item>
                <Descriptions.Item label="平台 Token">
                  {serverInfo.platformToken ? (
                    <Typography.Text copyable>{serverInfo.platformToken}</Typography.Text>
                  ) : (
                    '未设置'
                  )}
                </Descriptions.Item>
                <Descriptions.Item label="服务 Token">
                  {serverInfo.serviceToken ? (
                    <Typography.Text copyable>{serverInfo.serviceToken}</Typography.Text>
                  ) : (
                    '未设置'
                  )}
                </Descriptions.Item>
                <Descriptions.Item label="来源">
                  {SOURCE_TYPE_LABELS[serverInfo.source] ?? serverInfo.source}
                </Descriptions.Item>
                <Descriptions.Item label="URL 列表">
                  <Space direction="vertical" size={4}>
                    {(serverInfo.urls ?? []).map((url) => (
                      <Typography.Text key={url} copyable={{ text: url }}>
                        {url}
                      </Typography.Text>
                    ))}
                  </Space>
                </Descriptions.Item>
                <Descriptions.Item label="请求头">
                  <CodeBlock value={formatJson(serverInfo.headers)} />
                </Descriptions.Item>
                <Descriptions.Item label="安全配置">
                  <CodeBlock value={formatJson(serverInfo.security)} />
                </Descriptions.Item>
                <Descriptions.Item label="创建时间">{serverInfo.createdAt}</Descriptions.Item>
                <Descriptions.Item label="更新时间">{serverInfo.updatedAt}</Descriptions.Item>
              </Descriptions>

              <Card title="MCP 连接与 curl" size="small">
                <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                  <Alert
                    type="info"
                    showIcon
                    message="不要直接在浏览器里打开 MCP Endpoint"
                    description="浏览器直接访问会发起 GET，请求通常会返回 404。当前 MCP 网关需要 POST JSON-RPC，并带上 Accept: application/json, text/event-stream。"
                  />
                  <Alert
                    type="warning"
                    showIcon
                    message="UUID 和 connectToken 现在都可以作为 /gateway/:serverToken/mcp 的参数"
                    description="UUID 方式更适合本地调试、排障和稳定定位具体 Server；connectToken 方式更适合点击“生成 connectToken 连接信息”后发给客户端或外部接入方。"
                  />
                  <Alert
                    type="success"
                    showIcon
                    message="当前页面展示的是可直接复制的测试命令"
                    description="开发环境下如果前端运行在 3000 端口，页面里的 /gateway 地址会通过前端代理转发到后端。你也可以把域名替换成后端地址，例如本地默认的 http://localhost:9002。"
                  />

                  <Card title="方式一：使用 UUID 直接访问" size="small">
                    <Descriptions bordered size="small" column={1}>
                      <Descriptions.Item label="Server UUID">
                        <Typography.Text copyable>{serverInfo.uuid}</Typography.Text>
                      </Descriptions.Item>
                      <Descriptions.Item label="MCP Endpoint">
                        <Typography.Text copyable>{uuidEndpoint}</Typography.Text>
                      </Descriptions.Item>
                      <Descriptions.Item label="推荐场景">本地联调、排查某个具体 server、手工测试最方便。</Descriptions.Item>
                    </Descriptions>
                    <Space direction="vertical" size="middle" style={{ width: '100%', marginTop: 16 }}>
                      <div>
                        <Typography.Text strong>tools/list</Typography.Text>
                        <CodeBlock value={uuidToolsListCurl} />
                      </div>
                      <div>
                        <Typography.Text strong>tools/call</Typography.Text>
                        <Typography.Paragraph type="secondary" style={{ marginBottom: 8 }}>
                          默认示例会带上当前页面里第一个可见工具名。若该工具需要参数，请按下方 MCP 协议视图或工具详情里的 schema 补充 arguments。
                        </Typography.Paragraph>
                        <CodeBlock value={uuidToolsCallCurl} />
                      </div>
                    </Space>
                  </Card>

                  <Card title="方式二：使用 connectToken 访问" size="small">
                    {connectToken ? (
                      <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                        <Descriptions bordered size="small" column={1}>
                          <Descriptions.Item label="connectToken">
                            <Typography.Text copyable>{connectToken}</Typography.Text>
                          </Descriptions.Item>
                          <Descriptions.Item label="MCP Endpoint">
                            <Typography.Text copyable>{connectTokenEndpoint}</Typography.Text>
                          </Descriptions.Item>
                          <Descriptions.Item label="推荐场景">
                            适合把 token 化后的接入地址发给客户端或第三方。当前实现每次点击都会新建一条 connectToken 记录，建议以最新一次返回为准。
                          </Descriptions.Item>
                        </Descriptions>
                        <div>
                          <Typography.Text strong>tools/list</Typography.Text>
                          <CodeBlock value={connectTokenToolsListCurl} />
                        </div>
                        <div>
                          <Typography.Text strong>tools/call</Typography.Text>
                          <CodeBlock value={connectTokenToolsCallCurl} />
                        </div>
                      </Space>
                    ) : (
                      <Alert
                        type="info"
                        showIcon
                        message="还没有生成 connectToken"
                        description="点击上方“生成 connectToken 连接信息”后，这里会展示 token 化的 Endpoint 和对应的 curl 示例。UUID 方式现在已经可以直接使用。"
                      />
                    )}
                  </Card>
                </Space>
              </Card>
            </Space>
          ) : (
            !pageLoading && <Typography.Text type="secondary">暂无数据</Typography.Text>
          )}
        </Spin>
      ),
    },
    {
      key: 'tools',
      label: `工具管理 (${tools.length})`,
      children: (
        <Spin spinning={toolsLoading}>
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <Alert
              type="info"
              showIcon
              message="表单创建只负责 Server 信息"
              description="当前后端的表单接口只创建或更新 MCP Server；本页的“新增工具 / 编辑工具”用于维护 method、path、请求参数，以及接口鉴权/平台鉴权等工具级配置。"
            />
            <Space wrap>
              <Button type="primary" onClick={() => navigate(`/server/${uuid}/tool/create`)}>
                新增工具
              </Button>
              <Button onClick={openBatchAuth} disabled={!tools.length}>
                批量更新接口鉴权
              </Button>
              <Button onClick={() => void loadTools()}>刷新工具列表</Button>
            </Space>
            <Table<CommonToolItemInfo>
              rowKey="id"
              columns={toolColumns}
              dataSource={tools}
              scroll={{ x: 1200 }}
              pagination={{ pageSize: 10, showSizeChanger: true }}
            />
          </Space>
        </Spin>
      ),
    },
    {
      key: 'protocol',
      label: `MCP 协议视图 (${protocolTools.length})`,
      children: (
        <Spin spinning={protocolLoading}>
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <Alert
              type="success"
              showIcon
              message="这里展示的是最终注册到 MCP 协议层的工具 schema"
              description="可以先用本页提供的 tools/list / tools/call curl 校验网关链路，再到工具编辑页测试下游 REST 请求以及鉴权参数落点。"
            />
            <Button type="primary" onClick={() => void loadProtocolTools()}>
              刷新 MCP 协议视图
            </Button>
            {protocolTools.length ? (
              <Collapse items={protocolItems} />
            ) : (
              <Typography.Text type="secondary">当前没有可展示的 MCP 协议工具。</Typography.Text>
            )}
          </Space>
        </Spin>
      ),
    },
  ];

  return (
    <div style={{ padding: 24 }}>
      <Typography.Title level={4} style={{ marginTop: 0 }}>
        MCP Server 详情
      </Typography.Title>

      <Tabs items={tabItems} />

      <Modal
        title="编辑 MCP Server"
        open={editOpen}
        onCancel={() => setEditOpen(false)}
        onOk={() => void submitEdit()}
        confirmLoading={editSubmitting}
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true, message: '请输入名称' }]}>
            <Input />
          </Form.Item>
          <Form.Item
            name="description"
            label="描述"
            rules={[{ required: true, message: '请输入描述' }]}
          >
            <Input.TextArea rows={4} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="批量更新接口鉴权"
        open={batchAuthOpen}
        onCancel={() => setBatchAuthOpen(false)}
        onOk={() => void submitBatchAuth()}
        confirmLoading={batchAuthSubmitting}
        width={720}
        destroyOnClose
      >
        <Table<CommonToolItemInfo>
          rowKey="id"
          columns={batchAuthColumns}
          dataSource={tools}
          pagination={false}
          size="small"
        />
      </Modal>
    </div>
  );
}
