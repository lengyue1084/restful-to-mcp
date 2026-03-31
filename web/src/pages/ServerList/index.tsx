import { useCallback, useEffect, useMemo, useState } from 'react';
import {
  Alert,
  Button,
  Card,
  Descriptions,
  Input,
  message,
  Modal,
  Popconfirm,
  Space,
  Spin,
  Table,
  Tag,
  Typography,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useNavigate } from 'react-router-dom';
import {
  AUTH_TYPE_LABELS,
  SERVER_STATUS_COLORS,
  SERVER_STATUS_LABELS,
  SOURCE_TYPE_LABELS,
} from '../../types';
import type { GetMcpServerInfoByUUIDResponse, McpServerListItem } from '../../types';
import {
  deleteMcpServerByUUID,
  getMcpConnectTokenByUUID,
  getMcpServerInfoByUUID,
  listMcpServers,
} from '../../api';
import { buildMcpEndpoint } from '../../utils/mcp';

const RECENT_UUIDS_KEY = 'mcp-server-recent-uuids';

function loadRecentUuids(): string[] {
  try {
    const raw = localStorage.getItem(RECENT_UUIDS_KEY);
    if (!raw) return [];
    const parsed = JSON.parse(raw) as unknown;
    return Array.isArray(parsed)
      ? parsed.filter((value): value is string => typeof value === 'string')
      : [];
  } catch {
    return [];
  }
}

function saveRecentUuids(uuids: string[]) {
  localStorage.setItem(RECENT_UUIDS_KEY, JSON.stringify(uuids));
}

export default function ServerList() {
  const navigate = useNavigate();
  const [uuidInput, setUuidInput] = useState('');
  const [recentUuids, setRecentUuids] = useState<string[]>([]);
  const [serverInfo, setServerInfo] = useState<GetMcpServerInfoByUUIDResponse | null>(null);
  const [serverList, setServerList] = useState<McpServerListItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [listLoading, setListLoading] = useState(false);
  const [connectOpen, setConnectOpen] = useState(false);
  const [connectUrl, setConnectUrl] = useState('');
  const [connectLoading, setConnectLoading] = useState(false);

  const loadServerList = useCallback(async () => {
    setListLoading(true);
    try {
      const res = await listMcpServers();
      setServerList(res.data.data.items ?? []);
    } catch (error) {
      message.error(error instanceof Error ? error.message : '获取 Server 列表失败');
    } finally {
      setListLoading(false);
    }
  }, []);

  useEffect(() => {
    setRecentUuids(loadRecentUuids());
    void loadServerList();
  }, [loadServerList]);

  const rememberUuid = useCallback((uuid: string) => {
    const trimmed = uuid.trim();
    if (!trimmed) return;
    const next = [trimmed, ...recentUuids.filter((item) => item !== trimmed)].slice(0, 10);
    setRecentUuids(next);
    saveRecentUuids(next);
  }, [recentUuids]);

  const handleSearch = async (value?: string) => {
    const uuid = (value ?? uuidInput).trim();
    if (!uuid) {
      message.warning('请输入 Server UUID');
      return;
    }

    setLoading(true);
    setServerInfo(null);
    try {
      const res = await getMcpServerInfoByUUID(uuid);
      setServerInfo(res.data.data);
      rememberUuid(uuid);
    } catch (error) {
      message.error(error instanceof Error ? error.message : '查询失败');
    } finally {
      setLoading(false);
    }
  };

  const openConnectModal = async () => {
    if (!serverInfo) return;

    setConnectLoading(true);
    try {
      const res = await getMcpConnectTokenByUUID(serverInfo.uuid);
      setConnectUrl(buildMcpEndpoint(res.data.data.connectToken));
      setConnectOpen(true);
    } catch (error) {
      message.error(error instanceof Error ? error.message : '获取连接地址失败');
    } finally {
      setConnectLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!serverInfo) return;

    try {
      await deleteMcpServerByUUID(serverInfo.uuid);
      message.success('已删除');
      const next = recentUuids.filter((item) => item !== serverInfo.uuid);
      setRecentUuids(next);
      saveRecentUuids(next);
      setServerInfo(null);
      await loadServerList();
    } catch (error) {
      message.error(error instanceof Error ? error.message : '删除失败');
    }
  };

  const columns: ColumnsType<McpServerListItem> = useMemo(
    () => [
      {
        title: '名称',
        dataIndex: 'name',
        key: 'name',
        render: (_, record) => (
          <Button type="link" onClick={() => navigate(`/server/${record.uuid}`)}>
            {record.name || record.uuid}
          </Button>
        ),
      },
      {
        title: 'UUID',
        dataIndex: 'uuid',
        key: 'uuid',
        ellipsis: true,
      },
      {
        title: '来源',
        dataIndex: 'source',
        key: 'source',
        width: 120,
        render: (value: number) => SOURCE_TYPE_LABELS[value] ?? value,
      },
      {
        title: '认证',
        dataIndex: 'isAuth',
        key: 'isAuth',
        width: 120,
        render: (value: number) => AUTH_TYPE_LABELS[value] ?? value,
      },
      {
        title: '状态',
        dataIndex: 'status',
        key: 'status',
        width: 120,
        render: (value: number) => (
          <Tag color={SERVER_STATUS_COLORS[value]}>{SERVER_STATUS_LABELS[value] ?? value}</Tag>
        ),
      },
      {
        title: '工具数',
        dataIndex: 'toolCount',
        key: 'toolCount',
        width: 90,
      },
      {
        title: '更新时间',
        dataIndex: 'updatedAt',
        key: 'updatedAt',
        width: 180,
      },
    ],
    [navigate],
  );

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <div>
        <Typography.Title level={4} style={{ marginTop: 0, marginBottom: 8 }}>
          MCP Server
        </Typography.Title>
        <Typography.Text type="secondary">
          已补上 Server 列表接口，首页现在既能按 UUID 精确查询，也能直接浏览当前服务。
        </Typography.Text>
      </div>

      <Alert
        type="info"
        showIcon
        message="使用方式"
        description="列表适合浏览已有服务，UUID 查询适合快速定位单个服务。查询结果会保留最近记录，便于重复管理。"
      />

      <Input.Search
        placeholder="输入 Server UUID"
        value={uuidInput}
        onChange={(event) => setUuidInput(event.target.value)}
        onSearch={handleSearch}
        enterButton="查询"
        loading={loading}
        style={{ maxWidth: 640 }}
      />

      {recentUuids.length > 0 ? (
        <div>
          <Typography.Text type="secondary" style={{ marginRight: 8 }}>
            最近查询:
          </Typography.Text>
          <Space size={[8, 8]} wrap>
            {recentUuids.map((uuid) => (
              <Tag
                key={uuid}
                style={{ cursor: 'pointer' }}
                onClick={() => {
                  setUuidInput(uuid);
                  void handleSearch(uuid);
                }}
              >
                {uuid}
              </Tag>
            ))}
          </Space>
        </div>
      ) : null}

      <Card
        title="Server 列表"
        extra={(
          <Button onClick={() => void loadServerList()} loading={listLoading}>
            刷新列表
          </Button>
        )}
      >
        <Table<McpServerListItem>
          rowKey="uuid"
          loading={listLoading}
          columns={columns}
          dataSource={serverList}
          pagination={{ pageSize: 10, showSizeChanger: true }}
          scroll={{ x: 960 }}
        />
      </Card>

      <Spin spinning={loading}>
        {serverInfo ? (
          <Card
            title="查询结果"
            extra={(
              <Space wrap>
                <Button type="link" onClick={() => navigate(`/server/${serverInfo.uuid}`)}>
                  查看详情
                </Button>
                <Button type="link" loading={connectLoading} onClick={() => void openConnectModal()}>
                  获取连接地址
                </Button>
                {serverInfo.source === 2 ? (
                  <Button type="link" onClick={() => navigate(`/create/form/${serverInfo.uuid}`)}>
                    编辑完整配置
                  </Button>
                ) : null}
                <Popconfirm title="确定删除这个 MCP Server 吗？" onConfirm={() => void handleDelete()}>
                  <Button type="link" danger>
                    删除
                  </Button>
                </Popconfirm>
              </Space>
            )}
          >
            <Descriptions column={1} bordered size="small">
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
              <Descriptions.Item label="来源">
                {SOURCE_TYPE_LABELS[serverInfo.source] ?? serverInfo.source}
              </Descriptions.Item>
              <Descriptions.Item label="URL 数量">{serverInfo.urls?.length ?? 0}</Descriptions.Item>
              <Descriptions.Item label="创建时间">{serverInfo.createdAt}</Descriptions.Item>
              <Descriptions.Item label="更新时间">{serverInfo.updatedAt}</Descriptions.Item>
            </Descriptions>
          </Card>
        ) : null}
      </Spin>

      <Modal
        title="MCP 连接地址"
        open={connectOpen}
        onCancel={() => setConnectOpen(false)}
        footer={[
          <Button key="close" onClick={() => setConnectOpen(false)}>
            关闭
          </Button>,
        ]}
        destroyOnHidden
      >
        <Typography.Paragraph copyable={{ text: connectUrl }}>
          {connectUrl || '暂无'}
        </Typography.Paragraph>
      </Modal>
    </Space>
  );
}
