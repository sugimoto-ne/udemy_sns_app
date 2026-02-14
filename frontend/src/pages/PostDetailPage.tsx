import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Box, Typography, Divider, CircularProgress, Alert } from '@mui/material';
import { MainLayout } from '../components/layout/MainLayout';
import { PostCard } from '../components/post/PostCard';
import { CommentForm } from '../components/comment/CommentForm';
import { CommentList } from '../components/comment/CommentList';
import {
  usePost,
  useLikePost,
  useUnlikePost,
  useDeletePost,
} from '../hooks/usePosts';
import {
  useComments,
  useCreateComment,
  useDeleteComment,
} from '../hooks/useComments';
import type { CreateCommentRequest } from '../types/comment';

export const PostDetailPage: React.FC = () => {
  const { postId } = useParams<{ postId: string }>();
  const navigate = useNavigate();
  const postIdNum = postId ? parseInt(postId, 10) : 0;

  const { data: post, isLoading: postLoading, error: postError } = usePost(postIdNum);
  const { data: commentsData, isLoading: commentsLoading, error: commentsError } = useComments(postIdNum);

  const likePost = useLikePost();
  const unlikePost = useUnlikePost();
  const deletePost = useDeletePost();
  const createComment = useCreateComment();
  const deleteComment = useDeleteComment();

  const handleLike = (postId: number) => {
    likePost.mutate(postId);
  };

  const handleUnlike = (postId: number) => {
    unlikePost.mutate(postId);
  };

  const handleDeletePost = (postId: number) => {
    deletePost.mutate(postId, {
      onSuccess: () => {
        navigate('/');
      },
    });
  };

  const handleCreateComment = async (data: CreateCommentRequest) => {
    await createComment.mutateAsync({ postId: postIdNum, data });
  };

  const handleDeleteComment = (commentId: number) => {
    deleteComment.mutate({ postId: postIdNum, commentId });
  };

  if (postLoading) {
    return (
      <MainLayout>
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
          <CircularProgress />
        </Box>
      </MainLayout>
    );
  }

  if (postError || !post) {
    return (
      <MainLayout>
        <Alert severity="error">
          投稿の取得に失敗しました
        </Alert>
      </MainLayout>
    );
  }

  return (
    <MainLayout>
      <Box>
        {/* 投稿詳細 */}
        <PostCard
          post={post}
          onLike={handleLike}
          onUnlike={handleUnlike}
          onDelete={handleDeletePost}
        />

        <Divider sx={{ my: 3 }} />

        {/* コメントフォーム */}
        <Typography variant="h6" gutterBottom>
          コメント
        </Typography>
        <CommentForm
          onSubmit={handleCreateComment}
          isLoading={createComment.isPending}
        />

        {/* コメント一覧 */}
        <CommentList
          comments={commentsData?.data || []}
          isLoading={commentsLoading}
          error={commentsError}
          onDelete={handleDeleteComment}
        />
      </Box>
    </MainLayout>
  );
};
