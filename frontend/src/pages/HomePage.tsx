import React from 'react';
import { Box, CircularProgress, Alert, Typography, useTheme, useMediaQuery } from '@mui/material';
import { MainLayout } from '../components/layout/MainLayout';
import { PostForm } from '../components/post/PostForm';
import { PostCard } from '../components/post/PostCard';
import {
  useTimeline,
  useCreatePost,
  useLikePost,
  useUnlikePost,
  useDeletePost,
} from '../hooks/usePosts';
import type { CreatePostRequest } from '../types/post';

export const HomePage: React.FC = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const { data: timeline, isLoading, error } = useTimeline();
  const createPost = useCreatePost();
  const likePost = useLikePost();
  const unlikePost = useUnlikePost();
  const deletePost = useDeletePost();

  const handleCreatePost = async (data: CreatePostRequest) => {
    await createPost.mutateAsync(data);
  };

  const handleLike = (postId: number) => {
    likePost.mutate(postId);
  };

  const handleUnlike = (postId: number) => {
    unlikePost.mutate(postId);
  };

  const handleDelete = (postId: number) => {
    deletePost.mutate(postId);
  };

  return (
    <MainLayout>
      <Box>
        {/* 投稿フォーム（PC版のみ） */}
        {!isMobile && (
          <PostForm
            onSubmit={handleCreatePost}
            isLoading={createPost.isPending}
          />
        )}

        {/* タイムライン */}
        {isLoading && (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
            <CircularProgress />
          </Box>
        )}

        {error && (
          <Alert severity="error">
            タイムラインの取得に失敗しました
          </Alert>
        )}

        {timeline?.data && timeline.data.length === 0 && (
          <Box sx={{ textAlign: 'center', py: 4 }}>
            <Typography variant="body1" color="text.secondary">
              まだ投稿がありません
            </Typography>
          </Box>
        )}

        {timeline?.data && timeline.data.length > 0 && (
          <Box>
            {timeline.data.map((post) => (
              <PostCard
                key={post.id}
                post={post}
                onLike={handleLike}
                onUnlike={handleUnlike}
                onDelete={handleDelete}
              />
            ))}
          </Box>
        )}
      </Box>
    </MainLayout>
  );
};
