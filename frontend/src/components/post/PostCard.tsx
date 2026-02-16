import React from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Card,
  CardHeader,
  CardContent,
  CardActions,
  Avatar,
  Typography,
  IconButton,
  Box,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import {
  Favorite as FavoriteIcon,
  FavoriteBorder as FavoriteBorderIcon,
  ChatBubbleOutline as CommentIcon,
  Delete as DeleteIcon,
} from '@mui/icons-material';
import type { Post } from '../../types/post';
import { useAuth } from '../../contexts/AuthContext';
import { formatDistanceToNow } from 'date-fns';
import { ja } from 'date-fns/locale';
import { BookmarkButton } from './BookmarkButton';

interface PostCardProps {
  post: Post;
  onLike?: (postId: number) => void;
  onUnlike?: (postId: number) => void;
  onDelete?: (postId: number) => void;
}

export const PostCard: React.FC<PostCardProps> = ({
  post,
  onLike,
  onUnlike,
  onDelete,
}) => {
  const navigate = useNavigate();
  const { user: currentUser } = useAuth();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  const isOwner = currentUser?.id === post.user.id;

  const handleLikeClick = () => {
    if (post.is_liked && onUnlike) {
      onUnlike(post.id);
    } else if (!post.is_liked && onLike) {
      onLike(post.id);
    }
  };

  const handleDeleteClick = () => {
    if (onDelete && window.confirm('この投稿を削除しますか？')) {
      onDelete(post.id);
    }
  };

  const handleUserClick = () => {
    navigate(`/users/${post.user.username}`);
  };

  const handlePostClick = () => {
    navigate(`/posts/${post.id}`);
  };

  const formatDate = (dateString: string) => {
    try {
      return formatDistanceToNow(new Date(dateString), {
        addSuffix: true,
        locale: ja,
      });
    } catch {
      return dateString;
    }
  };

  return (
    <Card sx={{ mb: { xs: 1, sm: 2 } }} data-testid={`post-${post.id}`}>
      <CardHeader
        avatar={
          <Avatar
            src={post.user.avatar_url || undefined}
            alt={post.user.username}
            onClick={handleUserClick}
            sx={{
              cursor: 'pointer',
              width: { xs: 36, sm: 40 },
              height: { xs: 36, sm: 40 },
            }}
          >
            {post.user.username.charAt(0).toUpperCase()}
          </Avatar>
        }
        action={
          isOwner && (
            <IconButton onClick={handleDeleteClick} size={isMobile ? 'small' : 'medium'} data-testid="delete-post-button">
              <DeleteIcon />
            </IconButton>
          )
        }
        title={
          <Typography
            variant={isMobile ? 'body1' : 'subtitle1'}
            sx={{ cursor: 'pointer', fontWeight: 600 }}
            onClick={handleUserClick}
          >
            {post.user.display_name || post.user.username}
          </Typography>
        }
        subheader={
          <Typography variant="caption" color="text.secondary" sx={{ fontSize: { xs: '0.7rem', sm: '0.75rem' } }}>
            @{post.user.username} · {formatDate(post.created_at)}
          </Typography>
        }
        sx={{ pb: { xs: 1, sm: 2 } }}
      />

      <CardContent
        sx={{
          cursor: 'pointer',
          pt: 0,
          pb: { xs: 1, sm: 2 },
          px: { xs: 2, sm: 3 },
        }}
        onClick={handlePostClick}
      >
        <Typography
          variant="body1"
          sx={{
            whiteSpace: 'pre-wrap',
            fontSize: { xs: '0.9rem', sm: '1rem' },
            lineHeight: 1.5,
          }}
        >
          {post.content}
        </Typography>

        {/* メディア表示 */}
        {post.media && post.media.length > 0 && (
          <Box sx={{ mt: { xs: 1.5, sm: 2 }, display: 'flex', flexWrap: 'wrap', gap: 1 }}>
            {post.media.map((media) => (
              <Box key={media.id} sx={{ width: '100%' }}>
                {media.media_type === 'image' && (
                  <img
                    src={media.media_url}
                    alt="Post media"
                    style={{
                      width: '100%',
                      maxHeight: isMobile ? '300px' : '400px',
                      borderRadius: '8px',
                      objectFit: 'cover',
                    }}
                  />
                )}
                {media.media_type === 'video' && (
                  <video
                    src={media.media_url}
                    controls
                    style={{
                      width: '100%',
                      maxHeight: isMobile ? '300px' : '400px',
                      borderRadius: '8px',
                    }}
                  />
                )}
                {media.media_type === 'audio' && (
                  <audio src={media.media_url} controls style={{ width: '100%' }} />
                )}
              </Box>
            ))}
          </Box>
        )}
      </CardContent>

      <CardActions disableSpacing sx={{ px: { xs: 1, sm: 2 }, py: { xs: 0.5, sm: 1 } }}>
        <IconButton
          onClick={handleLikeClick}
          color={post.is_liked ? 'error' : 'default'}
          size={isMobile ? 'small' : 'medium'}
          data-testid={post.is_liked ? 'unlike-button' : 'like-button'}
        >
          {post.is_liked ? <FavoriteIcon fontSize={isMobile ? 'small' : 'medium'} /> : <FavoriteBorderIcon fontSize={isMobile ? 'small' : 'medium'} />}
        </IconButton>
        <Typography
          variant="body2"
          color="text.secondary"
          sx={{ mr: { xs: 1.5, sm: 2 }, fontSize: { xs: '0.8rem', sm: '0.875rem' } }}
          data-testid="like-count"
        >
          {post.likes_count}
        </Typography>

        <IconButton onClick={handlePostClick} size={isMobile ? 'small' : 'medium'} data-testid="comment-button">
          <CommentIcon fontSize={isMobile ? 'small' : 'medium'} />
        </IconButton>
        <Typography
          variant="body2"
          color="text.secondary"
          sx={{ fontSize: { xs: '0.8rem', sm: '0.875rem' } }}
          data-testid="comment-count"
        >
          {post.comments_count}
        </Typography>

        {/* ブックマークボタン */}
        <Box sx={{ marginLeft: 'auto' }}>
          <BookmarkButton postId={post.id} isBookmarked={post.is_bookmarked || false} />
        </Box>
      </CardActions>
    </Card>
  );
};
