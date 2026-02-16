import React from 'react';
import { useParams } from 'react-router-dom';
import { Container, Typography, Box, CircularProgress } from '@mui/material';
import TagIcon from '@mui/icons-material/Tag';
import { useHashtagPosts } from '../hooks/useHashtags';
import { useInView } from 'react-intersection-observer';

export const HashtagPage: React.FC = () => {
  const { hashtagName } = useParams<{ hashtagName: string }>();
  const decodedHashtagName = decodeURIComponent(hashtagName || '');

  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading, error } =
    useHashtagPosts(decodedHashtagName);

  const { ref, inView } = useInView();

  // 無限スクロール
  React.useEffect(() => {
    if (inView && hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  }, [inView, hasNextPage, isFetchingNextPage, fetchNextPage]);

  if (isLoading) {
    return (
      <Container maxWidth="md" sx={{ py: 4 }}>
        <Box display="flex" justifyContent="center">
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  if (error) {
    return (
      <Container maxWidth="md" sx={{ py: 4 }}>
        <Typography color="error">エラーが発生しました</Typography>
      </Container>
    );
  }

  const allPosts = data?.pages.flatMap((page) => page.posts) || [];

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Box display="flex" alignItems="center" gap={1} mb={3}>
        <TagIcon color="primary" />
        <Typography variant="h4">#{decodedHashtagName}</Typography>
      </Box>

      {allPosts.length === 0 ? (
        <Typography color="text.secondary">
          このハッシュタグの投稿はまだありません
        </Typography>
      ) : (
        <>
          {allPosts.map((post: any) => (
            <Box key={post.id} mb={2}>
              {/* 既存のPostCardコンポーネントを使用 */}
              {/* <PostCard post={post} /> */}
              <Typography>Post ID: {post.id}</Typography>
            </Box>
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
    </Container>
  );
};
