import React from 'react';
import { Box, CircularProgress, Alert, Typography, useTheme, useMediaQuery } from '@mui/material';
import { useQueryClient } from '@tanstack/react-query';
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
import { useUploadMedia } from '../hooks/useMedia';
import type { CreatePostRequest } from '../types/post';

export const HomePage: React.FC = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const queryClient = useQueryClient();
  const { data: timeline, isLoading, error } = useTimeline();
  const createPost = useCreatePost();
  const uploadMedia = useUploadMedia();
  const likePost = useLikePost();
  const unlikePost = useUnlikePost();
  const deletePost = useDeletePost();

  const handleCreatePost = async (data: CreatePostRequest, files?: File[]) => {
    let createdPostId: number | null = null;

    try {
      // まず投稿を作成（タイムライン更新は手動で行う）
      const newPost = await createPost.mutateAsync(data);
      createdPostId = newPost.id;

      // 画像がある場合はアップロード
      if (files && files.length > 0 && newPost) {
        try {
          // 画像アップロード完了後にタイムライン更新
          await uploadMedia.mutateAsync({
            postId: newPost.id,
            files,
          });
        } catch (uploadError) {
          // 画像アップロード失敗時は投稿を削除
          console.error('Failed to upload media, deleting post:', uploadError);
          if (createdPostId) {
            await deletePost.mutateAsync(createdPostId);
          }
          throw new Error('画像のアップロードに失敗しました。投稿を取り消しました。');
        }
      } else {
        // 画像がない場合は即座にタイムライン更新
        queryClient.invalidateQueries({ queryKey: ['timeline'] });
      }
    } catch (error: any) {
      console.error('Failed to create post:', error);
      // エラーメッセージを整形
      const errorMessage = error.message || error.response?.data?.error?.message || '投稿に失敗しました';
      throw new Error(errorMessage);
    }
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
            isLoading={createPost.isPending || uploadMedia.isPending}
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
