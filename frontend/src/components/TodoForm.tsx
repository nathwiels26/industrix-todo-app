import React, { useEffect } from 'react';
import {
  Modal,
  Form,
  Input,
  Select,
  DatePicker,
} from 'antd';
import dayjs from 'dayjs';
import { Todo, CreateTodoRequest, UpdateTodoRequest, Priority } from '../types';
import { useTodo } from '../contexts/TodoContext';

interface TodoFormProps {
  visible: boolean;
  onClose: () => void;
  todo?: Todo | null;
}

const { TextArea } = Input;
const { Option } = Select;

const TodoForm: React.FC<TodoFormProps> = ({ visible, onClose, todo }) => {
  const [form] = Form.useForm();
  const { createTodo, updateTodo, categories } = useTodo();

  useEffect(() => {
    if (visible) {
      if (todo) {
        form.setFieldsValue({
          title: todo.title,
          description: todo.description,
          priority: todo.priority,
          category_id: todo.category_id,
          due_date: todo.due_date ? dayjs(todo.due_date) : null,
        });
      } else {
        form.resetFields();
      }
    }
  }, [visible, todo, form]);

  const handleSubmit = async (values: {
    title: string;
    description?: string;
    priority?: Priority;
    category_id?: number;
    due_date?: dayjs.Dayjs;
  }) => {
    try {
      const data: CreateTodoRequest | UpdateTodoRequest = {
        title: values.title,
        description: values.description || '',
        priority: values.priority,
        category_id: values.category_id,
        due_date: values.due_date ? values.due_date.toISOString() : undefined,
      };

      if (todo) {
        await updateTodo(todo.id, data as UpdateTodoRequest);
      } else {
        await createTodo(data as CreateTodoRequest);
      }

      onClose();
      form.resetFields();
    } catch {
      // Error handled in context
    }
  };

  return (
    <Modal
      title={todo ? 'Edit Todo' : 'Create Todo'}
      open={visible}
      onCancel={onClose}
      onOk={() => form.submit()}
      okText={todo ? 'Update' : 'Create'}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        initialValues={{
          priority: 'medium',
        }}
      >
        <Form.Item
          name="title"
          label="Title"
          rules={[
            { required: true, message: 'Please enter a title' },
            { max: 255, message: 'Title must be less than 255 characters' },
          ]}
        >
          <Input placeholder="Enter todo title" />
        </Form.Item>

        <Form.Item
          name="description"
          label="Description"
        >
          <TextArea rows={3} placeholder="Enter description (optional)" />
        </Form.Item>

        <Form.Item
          name="priority"
          label="Priority"
        >
          <Select placeholder="Select priority">
            <Option value="high">High</Option>
            <Option value="medium">Medium</Option>
            <Option value="low">Low</Option>
          </Select>
        </Form.Item>

        <Form.Item
          name="category_id"
          label="Category"
        >
          <Select placeholder="Select category" allowClear>
            {categories.map((cat) => (
              <Option key={cat.id} value={cat.id}>
                <span style={{ color: cat.color }}>{cat.name}</span>
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          name="due_date"
          label="Due Date"
        >
          <DatePicker
            showTime
            style={{ width: '100%' }}
            placeholder="Select due date"
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default TodoForm;
