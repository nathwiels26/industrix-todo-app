import React from 'react';
import { ConfigProvider, Layout, Typography } from 'antd';
import { TodoProvider } from './contexts/TodoContext';
import TodoList from './components/TodoList';

const { Header, Content, Footer } = Layout;
const { Title } = Typography;

const App: React.FC = () => {
  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: '#3B82F6',
        },
      }}
    >
      <TodoProvider>
        <Layout style={{ minHeight: '100vh' }}>
          <Header
            style={{
              background: '#fff',
              padding: '0 24px',
              display: 'flex',
              alignItems: 'center',
              boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
            }}
          >
            <Title level={3} style={{ margin: 0, color: '#3B82F6' }}>
              Industrix Todo App
            </Title>
          </Header>
          <Content
            style={{
              padding: '24px',
              background: '#f0f2f5',
            }}
          >
            <div style={{ maxWidth: 1200, margin: '0 auto' }}>
              <TodoList />
            </div>
          </Content>
          <Footer style={{ textAlign: 'center', background: '#fff' }}>
            Industrix Todo App - Full Stack Challenge
          </Footer>
        </Layout>
      </TodoProvider>
    </ConfigProvider>
  );
};

export default App;
