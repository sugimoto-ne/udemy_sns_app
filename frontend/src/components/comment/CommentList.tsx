import React from 'react';
import {
  Box,
  Card,
  CardHeader,
  CardContent,
  Avatar,
  Typography,
  IconButton,
  CircularProgress,
  Alert,
} from '@mui/material';
import { Delete as DeleteIcon } from '@mui/icons-material';
import type { Comment } from '../../types/comment';
import { useAuth } from '../../contexts/AuthContext';
import { formatDistanceToNow } from 'date-fns';
import { ja } from 'date-fns/locale';

interface CommentListProps {
  comments: Comment[];
  isLoading?: boolean;
  error?: any;
  onDelete?: (commentId: number) => void;
}

export const CommentList: React.FC<CommentListProps> = ({
  comments,
  isLoading,
  error,
  onDelete,
}) => {
  const { user: currentUser } = useAuth();

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

  if (isLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 3 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        コメントの取得に失敗しました
      </Alert>
    );
  }

  if (comments.length === 0) {
    return (
      <Box sx={{ textAlign: 'center', py: 3 }}>
        <Typography variant="body2" color="text.secondary">
          まだコメントがありません
        </Typography>
      </Box>
    );
  }

  return (
    <Box>
      {comments.map((comment) => {
        const isOwner = currentUser?.id === comment.user.id;

        return (
          <Card key={comment.id} variant="outlined" sx={{ mb: 1.5 }}>
            <CardHeader
              avatar={
                <Avatar
                  src={comment.user.avatar_url || undefined}
                  alt={comment.user.username}
                  sx={{ width: 32, height: 32 }}
                >
                  {comment.user.username.charAt(0).toUpperCase()}
                </Avatar>
              }
              action={
                isOwner && onDelete && (
                  <IconButton
                    size="small"
                    onClick={() => {
                      if (window.confirm('このコメントを削除しますか？')) {
                        onDelete(comment.id);
                      }
                    }}
                  >
                    <DeleteIcon fontSize="small" />
                  </IconButton>
                )
              }
              title={
                <Typography variant="body2" fontWeight="bold">
                  {comment.user.display_name || comment.user.username}
                </Typography>
              }
              subheader={
                <Typography variant="caption" color="text.secondary">
                  @{comment.user.username} · {formatDate(comment.created_at)}
                </Typography>
              }
              sx={{ pb: 0 }}
            />
            <CardContent sx={{ pt: 1 }}>
              <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap' }}>
                {comment.content}
              </Typography>
            </CardContent>
          </Card>
        );
      })}
    </Box>
  );
};
