import React from 'react';
import { IconButton, Tooltip } from '@mui/material';
import BookmarkBorderIcon from '@mui/icons-material/BookmarkBorder';
import BookmarkIcon from '@mui/icons-material/Bookmark';
import { useBookmarkPost, useUnbookmarkPost } from '../../hooks/useBookmarks';

interface BookmarkButtonProps {
  postId: number;
  isBookmarked: boolean;
}

export const BookmarkButton: React.FC<BookmarkButtonProps> = ({ postId, isBookmarked }) => {
  const bookmarkMutation = useBookmarkPost();
  const unbookmarkMutation = useUnbookmarkPost();

  const handleClick = async (e: React.MouseEvent) => {
    e.stopPropagation(); // 親要素のクリックイベントを防ぐ

    if (isBookmarked) {
      await unbookmarkMutation.mutateAsync(postId);
    } else {
      await bookmarkMutation.mutateAsync(postId);
    }
  };

  const isLoading = bookmarkMutation.isPending || unbookmarkMutation.isPending;

  return (
    <Tooltip title={isBookmarked ? 'ブックマーク解除' : 'ブックマーク'}>
      <IconButton
        onClick={handleClick}
        disabled={isLoading}
        size="small"
        color={isBookmarked ? 'primary' : 'default'}
      >
        {isBookmarked ? <BookmarkIcon /> : <BookmarkBorderIcon />}
      </IconButton>
    </Tooltip>
  );
};
