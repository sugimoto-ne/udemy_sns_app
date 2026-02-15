import React from 'react';
import { Button } from '@mui/material';
import { useFollowUser, useUnfollowUser } from '../../hooks/useUsers';

interface FollowButtonProps {
  username: string;
  isFollowing: boolean;
}

export const FollowButton: React.FC<FollowButtonProps> = ({
  username,
  isFollowing,
}) => {
  const followUser = useFollowUser();
  const unfollowUser = useUnfollowUser();

  const handleClick = () => {
    if (isFollowing) {
      unfollowUser.mutate(username);
    } else {
      followUser.mutate(username);
    }
  };

  const isLoading = followUser.isPending || unfollowUser.isPending;

  return (
    <Button
      variant={isFollowing ? 'outlined' : 'contained'}
      size="small"
      onClick={handleClick}
      disabled={isLoading}
      data-testid={isFollowing ? 'unfollow-button' : 'follow-button'}
    >
      {isFollowing ? 'フォロー中' : 'フォロー'}
    </Button>
  );
};
