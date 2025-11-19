import React, { useState } from 'react';
import {
  Modal,
  Form,
  Input,
  Button,
  List,
  Popconfirm,
  ColorPicker,
  Space,
  Tag,
} from 'antd';
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { Category, CreateCategoryRequest, UpdateCategoryRequest } from '../types';
import { useTodo } from '../contexts/TodoContext';
import type { Color } from 'antd/es/color-picker';

interface CategoryManagerProps {
  visible: boolean;
  onClose: () => void;
}

const CategoryManager: React.FC<CategoryManagerProps> = ({ visible, onClose }) => {
  const [form] = Form.useForm();
  const { categories, createCategory, updateCategory, deleteCategory } = useTodo();
  const [editingId, setEditingId] = useState<number | null>(null);
  const [color, setColor] = useState<string>('#3B82F6');

  const handleSubmit = async (values: { name: string }) => {
    try {
      const data: CreateCategoryRequest | UpdateCategoryRequest = {
        name: values.name,
        color: color,
      };

      if (editingId) {
        await updateCategory(editingId, data as UpdateCategoryRequest);
        setEditingId(null);
      } else {
        await createCategory(data as CreateCategoryRequest);
      }

      form.resetFields();
      setColor('#3B82F6');
    } catch {
      // Error handled in context
    }
  };

  const handleEdit = (category: Category) => {
    setEditingId(category.id);
    setColor(category.color);
    form.setFieldsValue({ name: category.name });
  };

  const handleCancelEdit = () => {
    setEditingId(null);
    form.resetFields();
    setColor('#3B82F6');
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteCategory(id);
    } catch {
      // Error handled in context
    }
  };

  const handleColorChange = (value: Color) => {
    setColor(value.toHexString());
  };

  return (
    <Modal
      title="Manage Categories"
      open={visible}
      onCancel={onClose}
      footer={null}
      width={500}
    >
      <Form
        form={form}
        layout="inline"
        onFinish={handleSubmit}
        style={{ marginBottom: 16 }}
      >
        <Form.Item
          name="name"
          rules={[{ required: true, message: 'Please enter category name' }]}
          style={{ flex: 1 }}
        >
          <Input placeholder="Category name" />
        </Form.Item>
        <Form.Item>
          <ColorPicker value={color} onChange={handleColorChange} />
        </Form.Item>
        <Form.Item>
          <Space>
            <Button type="primary" htmlType="submit" icon={<PlusOutlined />}>
              {editingId ? 'Update' : 'Add'}
            </Button>
            {editingId && (
              <Button onClick={handleCancelEdit}>Cancel</Button>
            )}
          </Space>
        </Form.Item>
      </Form>

      <List
        dataSource={categories}
        renderItem={(category) => (
          <List.Item
            actions={[
              <Button
                key="edit"
                type="text"
                icon={<EditOutlined />}
                onClick={() => handleEdit(category)}
              />,
              <Popconfirm
                key="delete"
                title="Delete this category?"
                onConfirm={() => handleDelete(category.id)}
                okText="Yes"
                cancelText="No"
              >
                <Button type="text" danger icon={<DeleteOutlined />} />
              </Popconfirm>,
            ]}
          >
            <Tag color={category.color}>{category.name}</Tag>
          </List.Item>
        )}
      />
    </Modal>
  );
};

export default CategoryManager;
