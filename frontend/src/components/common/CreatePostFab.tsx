import React, { useState } from 'react';
import { Fab, useTheme, useMediaQuery } from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import { CreatePostModal } from '../post/CreatePostModal';
import { useCreatePost } from '../../hooks/usePosts';

export const CreatePostFab: React.FC = () => {
  const [open, setOpen] = useState(false);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const { mutateAsync: createPost, isPending } = useCreatePost();

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleSubmit = async (data: { content: string }) => {
    await createPost(data);
  };

  // モバイル時のみ表示
  if (!isMobile) {
    return null;
  }

  return (
    <>
      <Fab
        color="primary"
        aria-label="add"
        onClick={handleOpen}
        sx={{
          position: 'fixed',
          bottom: 16,
          right: 16,
          zIndex: 1000,
        }}
      >
        <AddIcon />
      </Fab>

      <CreatePostModal
        open={open}
        onClose={handleClose}
        onSubmit={handleSubmit}
        isLoading={isPending}
      />
    </>
  );
};
