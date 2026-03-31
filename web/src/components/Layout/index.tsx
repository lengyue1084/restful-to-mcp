import React from 'react';
import { Layout, Menu, Typography } from 'antd';
import {
  CloudServerOutlined,
  FileAddOutlined,
  FormOutlined,
  HomeOutlined,
} from '@ant-design/icons';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';

const { Sider, Header, Content } = Layout;

const menuItems = [
  { key: '/', icon: <HomeOutlined />, label: 'MCP Server' },
  { key: '/create/file', icon: <FileAddOutlined />, label: 'OpenAPI导入' },
  { key: '/create/form', icon: <FormOutlined />, label: '表单创建' },
];

const AppLayout: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const selectedKey = menuItems.find(
    (item) => item.key !== '/' && location.pathname.startsWith(item.key)
  )?.key || '/';

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="light" width={220} style={{ borderRight: '1px solid #f0f0f0' }}>
        <div style={{ padding: '20px 16px', textAlign: 'center' }}>
          <CloudServerOutlined style={{ fontSize: 28, color: '#1677ff' }} />
          <Typography.Title level={5} style={{ margin: '8px 0 0' }}>
            RESTful-to-MCP
          </Typography.Title>
        </div>
        <Menu
          mode="inline"
          selectedKeys={[selectedKey]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', borderBottom: '1px solid #f0f0f0', display: 'flex', alignItems: 'center' }}>
          <Typography.Text type="secondary">
            RESTful API to MCP Server Management
          </Typography.Text>
        </Header>
        <Content style={{ margin: 24, padding: 24, background: '#fff', borderRadius: 8, minHeight: 360 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
};

export default AppLayout;
