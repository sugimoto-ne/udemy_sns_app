import React from 'react';
import { useParams } from 'react-router-dom';
import {
  Box,
  Avatar,
  Typography,
  Paper,
  CircularProgress,
  Alert,
  Chip,
  Link,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import { MainLayout } from '../components/layout/MainLayout';
import { PostCard } from '../components/post/PostCard';
import { FollowButton } from '../components/user/FollowButton';
import { useUserProfile } from '../hooks/useUsers';
import { useUserPosts, useLikePost, useUnlikePost, useDeletePost } from '../hooks/usePosts';
import { useAuth } from '../contexts/AuthContext';
import { format } from 'date-fns';

export const UserProfilePage: React.FC = () => {
  const { username } = useParams<{ username: string }>();
  const { user: currentUser } = useAuth();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  const { data: profile, isLoading: profileLoading, error: profileError } = useUserProfile(username || '');
  const { data: postsData, isLoading: postsLoading, error: postsError } = useUserPosts(username || '');

  const likePost = useLikePost();
  const unlikePost = useUnlikePost();
  const deletePost = useDeletePost();

  const handleLike = (postId: number) => {
    likePost.mutate(postId);
  };

  const handleUnlike = (postId: number) => {
    unlikePost.mutate(postId);
  };

  const handleDelete = (postId: number) => {
    deletePost.mutate(postId);
  };

  if (profileLoading) {
    return (
      <MainLayout>
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
          <CircularProgress />
        </Box>
      </MainLayout>
    );
  }

  if (profileError || !profile) {
    return (
      <MainLayout>
        <Alert severity="error">
          ユーザーの取得に失敗しました
        </Alert>
      </MainLayout>
    );
  }

  const isOwnProfile = currentUser?.id === profile.id;

  return (
    <MainLayout>
      <Box>
        {/* プロフィールヘッダー */}
        <Paper sx={{ p: { xs: 2, sm: 3 }, mb: { xs: 2, sm: 3 } }}>
          {/* ヘッダー画像 */}
          {profile.header_url && (
            <Box
              sx={{
                width: '100%',
                height: { xs: 120, sm: 200 },
                backgroundImage: `url(${profile.header_url})`,
                backgroundSize: 'cover',
                backgroundPosition: 'center',
                borderRadius: 1,
                mb: 2,
              }}
            />
          )}

          {/* プロフィール情報 */}
          <Box sx={{ display: 'flex', flexDirection: { xs: 'column', sm: 'row' }, alignItems: { xs: 'center', sm: 'flex-start' }, gap: 2 }}>
            <Avatar
              src={profile.avatar_url || undefined}
              alt={profile.username}
              sx={{ width: { xs: 80, sm: 100 }, height: { xs: 80, sm: 100 } }}
            >
              {profile.username.charAt(0).toUpperCase()}
            </Avatar>

            <Box sx={{ flexGrow: 1, width: { xs: '100%', sm: 'auto' } }}>
              <Box sx={{ display: 'flex', flexDirection: { xs: 'column', sm: 'row' }, justifyContent: 'space-between', alignItems: { xs: 'center', sm: 'center' }, mb: 1, gap: { xs: 1, sm: 0 } }}>
                <Box sx={{ textAlign: { xs: 'center', sm: 'left' } }}>
                  <Typography variant={isMobile ? 'h6' : 'h5'} fontWeight="bold">
                    {profile.display_name || profile.username}
                  </Typography>
                  <Typography variant="body2" color="text.secondary" sx={{ fontSize: { xs: '0.8rem', sm: '0.875rem' } }}>
                    @{profile.username}
                  </Typography>
                </Box>
                {!isOwnProfile && profile.is_following !== undefined && (
                  <FollowButton username={profile.username} isFollowing={profile.is_following} />
                )}
              </Box>

              {profile.bio && (
                <Typography
                  variant="body1"
                  sx={{
                    mb: 2,
                    whiteSpace: 'pre-wrap',
                    fontSize: { xs: '0.9rem', sm: '1rem' },
                    textAlign: { xs: 'center', sm: 'left' },
                  }}
                >
                  {profile.bio}
                </Typography>
              )}

              <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap', mb: 1, justifyContent: { xs: 'center', sm: 'flex-start' } }}>
                {profile.occupation && (
                  <Chip label={profile.occupation} size="small" />
                )}
                {profile.website && (
                  <Link href={profile.website} target="_blank" rel="noopener noreferrer" sx={{ fontSize: { xs: '0.8rem', sm: '0.875rem' } }}>
                    {profile.website}
                  </Link>
                )}
              </Box>

              <Box sx={{ display: 'flex', gap: { xs: 2, sm: 3 }, justifyContent: { xs: 'center', sm: 'flex-start' } }}>
                <Typography variant="body2" sx={{ fontSize: { xs: '0.8rem', sm: '0.875rem' } }}>
                  <strong>{profile.following_count}</strong> フォロー中
                </Typography>
                <Typography variant="body2" sx={{ fontSize: { xs: '0.8rem', sm: '0.875rem' } }}>
                  <strong>{profile.followers_count}</strong> フォロワー
                </Typography>
              </Box>

              <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 1, fontSize: { xs: '0.7rem', sm: '0.75rem' }, textAlign: { xs: 'center', sm: 'left' } }}>
                {profile.created_at && `登録日: ${format(new Date(profile.created_at), 'yyyy年M月d日')}`}
              </Typography>
            </Box>
          </Box>
        </Paper>

        {/* ユーザーの投稿一覧 */}
        <Typography variant="h6" gutterBottom sx={{ fontSize: { xs: '1.1rem', sm: '1.25rem' } }}>
          投稿
        </Typography>

        {postsLoading && (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
            <CircularProgress />
          </Box>
        )}

        {postsError && (
          <Alert severity="error">
            投稿の取得に失敗しました
          </Alert>
        )}

        {postsData?.data && postsData.data.length === 0 && (
          <Box sx={{ textAlign: 'center', py: 4 }}>
            <Typography variant="body1" color="text.secondary">
              まだ投稿がありません
            </Typography>
          </Box>
        )}

        {postsData?.data && postsData.data.length > 0 && (
          <Box>
            {postsData.data.map((post) => (
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
