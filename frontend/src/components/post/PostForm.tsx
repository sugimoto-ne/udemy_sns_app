import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import {
  Box,
  Button,
  TextField,
  Paper,
  Typography,
  Alert,
} from '@mui/material';
import type { CreatePostRequest } from '../../types/post';

interface PostFormProps {
  onSubmit: (data: CreatePostRequest) => Promise<void>;
  initialContent?: string;
  isLoading?: boolean;
}

export const PostForm: React.FC<PostFormProps> = ({
  onSubmit,
  initialContent = '',
  isLoading = false,
}) => {
  const [error, setError] = useState<string>('');

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CreatePostRequest>({
    defaultValues: {
      content: initialContent,
    },
  });

  const handleFormSubmit = async (data: CreatePostRequest) => {
    try {
      setError('');
      await onSubmit(data);
      reset();
    } catch (err: any) {
      setError(
        err.response?.data?.error?.message || '投稿に失敗しました'
      );
    }
  };

  return (
    <Paper sx={{ p: { xs: 2, sm: 3 }, mb: { xs: 2, sm: 3 } }}>
      <Typography variant="h6" gutterBottom sx={{ fontSize: { xs: '1.1rem', sm: '1.25rem' } }}>
        新規投稿
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <form onSubmit={handleSubmit(handleFormSubmit)}>
        <TextField
          fullWidth
          multiline
          rows={4}
          placeholder="いまどうしてる？"
          {...register('content', {
            required: '投稿内容を入力してください',
            maxLength: {
              value: 280,
              message: '投稿は280文字以内で入力してください',
            },
          })}
          error={!!errors.content}
          helperText={errors.content?.message}
          sx={{
            mb: 2,
            '& .MuiInputBase-input': {
              fontSize: { xs: '0.9rem', sm: '1rem' },
            },
          }}
        />

        <Box sx={{ display: 'flex', justifyContent: 'flex-end' }}>
          <Button
            type="submit"
            variant="contained"
            color="primary"
            disabled={isLoading}
            size="medium"
            sx={{ px: { xs: 2, sm: 3 } }}
          >
            {isLoading ? '投稿中...' : '投稿する'}
          </Button>
        </Box>
      </form>
    </Paper>
  );
};
