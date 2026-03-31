import { useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Alert, Button, Card, Form, Input, Select, Space, Typography, message } from 'antd';
import { v4 as uuidv4 } from 'uuid';
import {
  createMcpServerByForm,
  getMcpServerInfoByUUID,
  updateMcpServerByForm,
} from '../../api';
import type {
  AuthMode,
  AuthPosition,
  CreateMcpServerByFormRequest,
  Security,
} from '../../types';
import { AUTH_TYPE_LABELS } from '../../types';

const authOptions = Object.keys(AUTH_TYPE_LABELS).map((value) => {
  const authValue = Number(value);
  return {
    value: authValue,
    label: AUTH_TYPE_LABELS[authValue],
  };
});

function parseLines(value: string): string[] {
  return value
    .split(/\r?\n|,/)
    .map((item) => item.trim())
    .filter(Boolean);
}

function parseHeaders(raw: string): Record<string, string> {
  const text = raw.trim();
  if (!text) return {};
  const obj = JSON.parse(text) as Record<string, unknown>;
  const result: Record<string, string> = {};
  Object.keys(obj).forEach((key) => {
    result[key] = String(obj[key]);
  });
  return result;
}

export default function CreateByForm() {
  const { uuid } = useParams<{ uuid: string }>();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const isEdit = !!uuid;
  const isAuth = Form.useWatch('isAuth', form) as number | undefined;
  const authType = Number(isAuth);
  const securityMode = Form.useWatch('securityMode', form) as AuthMode | undefined;

  useEffect(() => {
    if (!uuid) return;
    setLoading(true);
    void (async () => {
      try {
        const res = await getMcpServerInfoByUUID(uuid);
        const data = res.data.data;
        form.setFieldsValue({
          uuid: data.uuid,
          name: data.name,
          description: data.description,
          urlsText: (data.urls ?? []).join('\n'),
          version: data.version || 'v1',
          isAuth: data.isAuth,
          platformToken: data.platformToken || '',
          serviceToken: data.serviceToken || '',
          headersText: JSON.stringify(data.headers ?? {}, null, 2),
          securityKey: data.security?.securityKey || '',
          securityMode: data.security?.mode || '',
          securityName: data.security?.name || '',
          securityIn: data.security?.in || 'header',
          securityScheme: data.security?.scheme || '',
        });
      } catch (e) {
        message.error(e instanceof Error ? e.message : '加载失败');
      } finally {
        setLoading(false);
      }
    })();
  }, [uuid, form]);

  const securityModeOptions = useMemo(
    () => [
      { label: 'apiKey', value: 'apiKey' },
      { label: 'http', value: 'http' },
    ],
    [],
  );

  const onFinish = async (values: any) => {
    setSubmitting(true);
    try {
      const payload: CreateMcpServerByFormRequest = {
        uuid: values.uuid || uuidv4(),
        name: values.name,
        description: values.description || '',
        urls: parseLines(values.urlsText || ''),
        version: values.version || 'v1',
        isAuth: Number(values.isAuth) as CreateMcpServerByFormRequest['isAuth'],
        platformToken: values.platformToken || '',
        serviceToken: values.serviceToken || '',
        headers: parseHeaders(values.headersText || ''),
      };

      if (Number(values.isAuth) === 3 || Number(values.isAuth) === 4) {
        const security: Security = {
          securityKey: values.securityKey || '',
          mode: values.securityMode || '',
          name: values.securityName || '',
          in: (values.securityIn || 'header') as AuthPosition,
          scheme: values.securityScheme || '',
          description: '',
          bearerFormat: '',
        };
        payload.security = security;
      }

      if (isEdit) {
        await updateMcpServerByForm(payload);
        message.success('更新成功');
      } else {
        await createMcpServerByForm(payload);
        message.success('创建成功，请继续新增工具并配置请求入参');
      }
      navigate(`/server/${payload.uuid}`);
    } catch (e) {
      message.error(e instanceof Error ? e.message : isEdit ? '更新失败' : '创建失败');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div>
      <Typography.Title level={4} style={{ marginTop: 0 }}>
        {isEdit ? '表单编辑 MCP Server' : '表单创建 MCP Server'}
      </Typography.Title>

      <Card loading={loading}>
        <Space direction="vertical" size="middle" style={{ width: '100%' }}>
          <Alert
            type="info"
            showIcon
            message="当前页面只负责创建 MCP Server"
            description="根据后端接口设计，表单创建接口不包含工具定义。创建完成后，请到详情页点击“新增工具”配置 method、path 和请求入参。后端当前没有单独保存工具出参结构的接口。"
          />

          <Form
            form={form}
            layout="vertical"
            initialValues={{
              uuid: uuidv4(),
              version: 'v1',
              isAuth: 1,
              headersText: '{\n  "Content-Type": "application/json"\n}',
              securityIn: 'header',
            }}
            onFinish={onFinish}
          >
            <Form.Item
              name="uuid"
              label="UUID"
              rules={[{ required: true, message: '请输入 UUID' }]}
            >
              <Input disabled={isEdit} />
            </Form.Item>

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
              name="urlsText"
              label="URL 列表"
              rules={[{ required: true, message: '请至少输入一个 URL' }]}
            >
              <Input.TextArea rows={4} placeholder="每行一个 URL，或使用英文逗号分隔" />
            </Form.Item>

            <Form.Item
              name="version"
              label="版本"
              rules={[{ required: true, message: '请输入版本' }]}
            >
              <Input />
            </Form.Item>

            <Form.Item
              name="isAuth"
              label="鉴权类型"
              rules={[{ required: true, message: '请选择鉴权类型' }]}
            >
              <Select options={authOptions} />
            </Form.Item>

            {(authType === 2 || authType === 4) && (
              <Form.Item
                name="platformToken"
                label="平台 Token"
                rules={[{ required: true, message: '请输入平台 Token' }]}
              >
                <Input />
              </Form.Item>
            )}

            {(authType === 3 || authType === 4) && (
              <>
                <Form.Item
                  name="serviceToken"
                  label="Service Token"
                  rules={[{ required: true, message: '请输入 Service Token' }]}
                >
                  <Input />
                </Form.Item>

                <Form.Item name="securityKey" label="Security Key">
                  <Input placeholder="例如 myAuth" />
                </Form.Item>

                <Form.Item
                  name="securityMode"
                  label="Security Mode"
                  rules={[{ required: true, message: '请选择 Security Mode' }]}
                >
                  <Select options={securityModeOptions} />
                </Form.Item>

                {securityMode === 'apiKey' ? (
                  <>
                    <Form.Item
                      name="securityName"
                      label="apiKey Name"
                      rules={[{ required: true, message: '请输入 apiKey Name' }]}
                    >
                      <Input />
                    </Form.Item>

                    <Form.Item
                      name="securityIn"
                      label="apiKey In"
                      rules={[{ required: true, message: '请选择位置' }]}
                    >
                      <Select
                        options={[
                          { label: 'header', value: 'header' },
                          { label: 'query', value: 'query' },
                          { label: 'cookie', value: 'cookie' },
                        ]}
                      />
                    </Form.Item>
                  </>
                ) : null}

                {securityMode === 'http' ? (
                  <Form.Item
                    name="securityScheme"
                    label="HTTP Scheme"
                    rules={[{ required: true, message: '请输入 scheme，例如 bearer' }]}
                  >
                    <Input />
                  </Form.Item>
                ) : null}
              </>
            )}

            <Form.Item
              name="headersText"
              label="Headers (JSON)"
              rules={[
                {
                  validator: async (_, value) => {
                    try {
                      parseHeaders(value || '');
                    } catch {
                      throw new Error('Headers 必须是合法的 JSON 对象');
                    }
                  },
                },
              ]}
            >
              <Input.TextArea rows={6} />
            </Form.Item>

            <Space>
              <Button type="primary" htmlType="submit" loading={submitting}>
                {isEdit ? '保存' : '创建'}
              </Button>
              <Button onClick={() => navigate(-1)}>返回</Button>
            </Space>
          </Form>
        </Space>
      </Card>
    </div>
  );
}
