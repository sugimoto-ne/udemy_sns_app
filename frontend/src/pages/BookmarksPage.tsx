import React from 'react';
import { Box, Typography, CircularProgress } from '@mui/material';
import BookmarkIcon from '@mui/icons-material/Bookmark';
import { useBookmarks } from '../hooks/useBookmarks';
import { useInView } from 'react-intersection-observer';
import { MainLayout } from '../components/layout/MainLayout';
import { PostCard } from '../components/post/PostCard';
import { useLikePost, useUnlikePost, useDeletePost } from '../hooks/usePosts';

export const BookmarksPage: React.FC = () => {
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading, error } =
    useBookmarks();
  const likePost = useLikePost();
  const unlikePost = useUnlikePost();
  const deletePost = useDeletePost();

  const { ref, inView } = useInView();

  // 無限スクロール
  React.useEffect(() => {
    if (inView && hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  }, [inView, hasNextPage, isFetchingNextPage, fetchNextPage]);

  const handleLike = (postId: number) => {
    likePost.mutate(postId);
  };

  const handleUnlike = (postId: number) => {
    unlikePost.mutate(postId);
  };

  const handleDelete = (postId: number) => {
    deletePost.mutate(postId);
  };

  if (isLoading) {
    return (
      <MainLayout>
        <Box display="flex" justifyContent="center" py={4}>
          <CircularProgress />
        </Box>
      </MainLayout>
    );
  }

  if (error) {
    return (
      <MainLayout>
        <Typography color="error">エラーが発生しました</Typography>
      </MainLayout>
    );
  }

  const allPosts = data?.pages.flatMap((page) => page.posts) || [];

  return (
    <MainLayout>
      <Box display="flex" alignItems="center" gap={1} mb={3}>
        <BookmarkIcon color="primary" />
        <Typography variant="h4">ブックマーク</Typography>
      </Box>

      {allPosts.length === 0 ? (
        <Typography color="text.secondary">ブックマークした投稿はありません</Typography>
      ) : (
        <>
          {allPosts.map((post) => (
            <PostCard
              key={post.id}
              post={post}
              onLike={handleLike}
              onUnlike={handleUnlike}
              onDelete={handleDelete}
            />
          ))}

          {/* 無限スクロールのトリガー */}
          <Box ref={ref} py={2} display="flex" justifyContent="center">
            {isFetchingNextPage && <CircularProgress />}
          </Box>

          {!hasNextPage && allPosts.length > 0 && (
            <Typography textAlign="center" color="text.secondary">
              全て表示しました
            </Typography>
          )}
        </>
      )}
    </MainLayout>
  );
};
