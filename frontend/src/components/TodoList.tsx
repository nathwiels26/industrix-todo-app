import React, { useEffect, useState } from 'react';
import {
  Table,
  Button,
  Tag,
  Space,
  Popconfirm,
  Input,
  Select,
  Row,
  Col,
  Card,
  Checkbox,
  Typography,
} from 'antd';
import {
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
  SearchOutlined,
  AppstoreOutlined,
} from '@ant-design/icons';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import dayjs from 'dayjs';
import { Todo, Priority } from '../types';
import { useTodo } from '../contexts/TodoContext';
import TodoForm from './TodoForm';
import CategoryManager from './CategoryManager';

const { Option } = Select;
const { Title } = Typography;

const priorityColors: Record<Priority, string> = {
  high: 'red',
  medium: 'orange',
  low: 'green',
};

const TodoList: React.FC = () => {
  const {
    todos,
    categories,
    pagination,
    loading,
    filter,
    fetchTodos,
    fetchCategories,
    deleteTodo,
    toggleComplete,
    setFilter,
  } = useTodo();

  const [formVisible, setFormVisible] = useState(false);
  const [categoryManagerVisible, setCategoryManagerVisible] = useState(false);
  const [editingTodo, setEditingTodo] = useState<Todo | null>(null);
  const [searchText, setSearchText] = useState('');

  useEffect(() => {
    fetchCategories();
    fetchTodos();
  }, [fetchCategories, fetchTodos]);

  const handleSearch = () => {
    const newFilter = { ...filter, search: searchText, page: 1 };
    setFilter(newFilter);
    fetchTodos(newFilter);
  };

  const handleTableChange = (paginationConfig: TablePaginationConfig) => {
    const newFilter = {
      ...filter,
      page: paginationConfig.current || 1,
      limit: paginationConfig.pageSize || 10,
    };
    setFilter(newFilter);
    fetchTodos(newFilter);
  };

  const handleCategoryFilter = (categoryId: number | undefined) => {
    const newFilter = { ...filter, category_id: categoryId, page: 1 };
    setFilter(newFilter);
    fetchTodos(newFilter);
  };

  const handleCompletedFilter = (completed: boolean | undefined) => {
    const newFilter = { ...filter, completed, page: 1 };
    setFilter(newFilter);
    fetchTodos(newFilter);
  };

  const handlePriorityFilter = (priority: Priority | undefined) => {
    const newFilter = { ...filter, priority, page: 1 };
    setFilter(newFilter);
    fetchTodos(newFilter);
  };

  const handleEdit = (todo: Todo) => {
    setEditingTodo(todo);
    setFormVisible(true);
  };

  const handleDelete = async (id: number) => {
    await deleteTodo(id);
  };

  const handleFormClose = () => {
    setFormVisible(false);
    setEditingTodo(null);
  };

  const columns: ColumnsType<Todo> = [
    {
      title: 'Status',
      key: 'completed',
      width: 80,
      render: (_, record) => (
        <Checkbox
          checked={record.completed}
          onChange={() => toggleComplete(record.id)}
        />
      ),
    },
    {
      title: 'Title',
      dataIndex: 'title',
      key: 'title',
      render: (text, record) => (
        <span style={{ textDecoration: record.completed ? 'line-through' : 'none' }}>
          {text}
        </span>
      ),
    },
    {
      title: 'Description',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      responsive: ['md'],
    },
    {
      title: 'Priority',
      dataIndex: 'priority',
      key: 'priority',
      width: 100,
      render: (priority: Priority) => (
        <Tag color={priorityColors[priority]}>{priority.toUpperCase()}</Tag>
      ),
    },
    {
      title: 'Category',
      key: 'category',
      width: 120,
      responsive: ['sm'],
      render: (_, record) =>
        record.category ? (
          <Tag color={record.category.color}>{record.category.name}</Tag>
        ) : (
          '-'
        ),
    },
    {
      title: 'Due Date',
      dataIndex: 'due_date',
      key: 'due_date',
      width: 120,
      responsive: ['lg'],
      render: (date) => (date ? dayjs(date).format('MMM D, YYYY') : '-'),
    },
    {
      title: 'Actions',
      key: 'actions',
      width: 100,
      render: (_, record) => (
        <Space>
          <Button
            type="text"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          />
          <Popconfirm
            title="Delete this todo?"
            onConfirm={() => handleDelete(record.id)}
            okText="Yes"
            cancelText="No"
          >
            <Button type="text" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card>
      <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
        <Col>
          <Title level={3} style={{ margin: 0 }}>Todo List</Title>
        </Col>
        <Col>
          <Space>
            <Button
              icon={<AppstoreOutlined />}
              onClick={() => setCategoryManagerVisible(true)}
            >
              Categories
            </Button>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => setFormVisible(true)}
            >
              Add Todo
            </Button>
          </Space>
        </Col>
      </Row>

      {/* Filters */}
      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} sm={12} md={6}>
          <Input
            placeholder="Search todos..."
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            onPressEnter={handleSearch}
            suffix={
              <SearchOutlined
                style={{ cursor: 'pointer' }}
                onClick={handleSearch}
              />
            }
          />
        </Col>
        <Col xs={12} sm={6} md={4}>
          <Select
            placeholder="Category"
            allowClear
            style={{ width: '100%' }}
            onChange={handleCategoryFilter}
          >
            {categories.map((cat) => (
              <Option key={cat.id} value={cat.id}>
                {cat.name}
              </Option>
            ))}
          </Select>
        </Col>
        <Col xs={12} sm={6} md={4}>
          <Select
            placeholder="Status"
            allowClear
            style={{ width: '100%' }}
            onChange={handleCompletedFilter}
          >
            <Option value={false}>Pending</Option>
            <Option value={true}>Completed</Option>
          </Select>
        </Col>
        <Col xs={12} sm={6} md={4}>
          <Select
            placeholder="Priority"
            allowClear
            style={{ width: '100%' }}
            onChange={handlePriorityFilter}
          >
            <Option value="high">High</Option>
            <Option value="medium">Medium</Option>
            <Option value="low">Low</Option>
          </Select>
        </Col>
      </Row>

      {/* Todo Table */}
      <Table
        columns={columns}
        dataSource={todos}
        rowKey="id"
        loading={loading}
        pagination={{
          current: pagination.current_page,
          pageSize: pagination.per_page,
          total: pagination.total,
          showSizeChanger: true,
          showTotal: (total) => `Total ${total} items`,
          pageSizeOptions: ['10', '20', '50'],
        }}
        onChange={handleTableChange}
        scroll={{ x: 800 }}
      />

      {/* Todo Form Modal */}
      <TodoForm
        visible={formVisible}
        onClose={handleFormClose}
        todo={editingTodo}
      />

      {/* Category Manager Modal */}
      <CategoryManager
        visible={categoryManagerVisible}
        onClose={() => setCategoryManagerVisible(false)}
      />
    </Card>
  );
};

export default TodoList;
