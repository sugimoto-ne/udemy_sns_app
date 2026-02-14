import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import {
  Box,
  Button,
  TextField,
  Alert,
} from '@mui/material';
import type { CreateCommentRequest } from '../../types/comment';

interface CommentFormProps {
  onSubmit: (data: CreateCommentRequest) => Promise<void>;
  isLoading?: boolean;
}

export const CommentForm: React.FC<CommentFormProps> = ({
  onSubmit,
  isLoading = false,
}) => {
  const [error, setError] = useState<string>('');

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CreateCommentRequest>();

  const handleFormSubmit = async (data: CreateCommentRequest) => {
    try {
      setError('');
      await onSubmit(data);
      reset();
    } catch (err: any) {
      setError(
        err.response?.data?.error?.message || 'コメントの投稿に失敗しました'
      );
    }
  };

  return (
    <Box sx={{ mb: 3 }}>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <form onSubmit={handleSubmit(handleFormSubmit)}>
        <TextField
          fullWidth
          multiline
          rows={2}
          placeholder="コメントを入力..."
          {...register('content', {
            required: 'コメントを入力してください',
            maxLength: {
              value: 280,
              message: 'コメントは280文字以内で入力してください',
            },
          })}
          error={!!errors.content}
          helperText={errors.content?.message}
          sx={{ mb: 1 }}
        />

        <Box sx={{ display: 'flex', justifyContent: 'flex-end' }}>
          <Button
            type="submit"
            variant="contained"
            size="small"
            disabled={isLoading}
          >
            {isLoading ? 'コメント中...' : 'コメント'}
          </Button>
        </Box>
      </form>
    </Box>
  );
};
