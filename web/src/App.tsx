import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import AppLayout from './components/Layout';
import ServerList from './pages/ServerList';
import CreateByFile from './pages/CreateByFile';
import CreateByForm from './pages/CreateByForm';
import ServerDetail from './pages/ServerDetail';
import ToolEdit from './pages/ToolEdit';

function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <BrowserRouter>
        <Routes>
          <Route element={<AppLayout />}>
            <Route path="/" element={<ServerList />} />
            <Route path="/create/file" element={<CreateByFile />} />
            <Route path="/create/form" element={<CreateByForm />} />
            <Route path="/create/form/:uuid" element={<CreateByForm />} />
            <Route path="/server/:uuid" element={<ServerDetail />} />
            <Route path="/server/:serverUUID/tool/create" element={<ToolEdit />} />
            <Route path="/server/:serverUUID/tool/:toolUUID" element={<ToolEdit />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  );
}

export default App;
