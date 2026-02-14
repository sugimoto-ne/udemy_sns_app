import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  TextField,
  Button,
  Box,
  Alert,
  IconButton,
  Slide,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import { Close as CloseIcon } from '@mui/icons-material';
import type { TransitionProps } from '@mui/material/transitions';
import type { CreatePostRequest } from '../../types/post';

interface CreatePostModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreatePostRequest) => Promise<void>;
  isLoading?: boolean;
}

const Transition = React.forwardRef(function Transition(
  props: TransitionProps & {
    children: React.ReactElement;
  },
  ref: React.Ref<unknown>,
) {
  return <Slide direction="up" ref={ref} {...props} />;
});

export const CreatePostModal: React.FC<CreatePostModalProps> = ({
  open,
  onClose,
  onSubmit,
  isLoading = false,
}) => {
  const [error, setError] = useState<string>('');
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CreatePostRequest>({
    defaultValues: {
      content: '',
    },
  });

  const handleFormSubmit = async (data: CreatePostRequest) => {
    try {
      setError('');
      await onSubmit(data);
      reset();
      onClose();
    } catch (err: any) {
      setError(
        err.response?.data?.error?.message || '投稿に失敗しました'
      );
    }
  };

  const handleClose = () => {
    reset();
    setError('');
    onClose();
  };

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      TransitionComponent={Transition}
      fullWidth
      maxWidth="sm"
      fullScreen={isMobile}
      PaperProps={{
        sx: {
          position: isMobile ? 'fixed' : 'relative',
          bottom: isMobile ? 0 : 'auto',
          m: isMobile ? 0 : 2,
          maxHeight: isMobile ? '80vh' : '90vh',
          borderRadius: isMobile ? '16px 16px 0 0' : '8px',
        },
      }}
    >
      <DialogTitle sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        新規投稿
        <IconButton
          edge="end"
          color="inherit"
          onClick={handleClose}
          aria-label="close"
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>

      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <form onSubmit={handleSubmit(handleFormSubmit)}>
          <TextField
            fullWidth
            multiline
            rows={6}
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
            autoFocus
            sx={{ mb: 3 }}
          />

          <Box sx={{ display: 'flex', justifyContent: 'flex-end', gap: 2 }}>
            <Button
              variant="outlined"
              onClick={handleClose}
              disabled={isLoading}
            >
              キャンセル
            </Button>
            <Button
              type="submit"
              variant="contained"
              color="primary"
              disabled={isLoading}
            >
              {isLoading ? '投稿中...' : '投稿する'}
            </Button>
          </Box>
        </form>
      </DialogContent>
    </Dialog>
  );
};
