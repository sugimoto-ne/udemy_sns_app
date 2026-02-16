import React from 'react';
import { Chip } from '@mui/material';
import { useNavigate } from 'react-router-dom';

interface HashtagChipProps {
  hashtag: string;
  onClick?: () => void;
}

export const HashtagChip: React.FC<HashtagChipProps> = ({ hashtag, onClick }) => {
  const navigate = useNavigate();

  const handleClick = () => {
    if (onClick) {
      onClick();
    } else {
      // ハッシュタグページに遷移
      navigate(`/hashtags/${encodeURIComponent(hashtag)}`);
    }
  };

  return (
    <Chip
      label={`#${hashtag}`}
      onClick={handleClick}
      size="small"
      sx={{
        cursor: 'pointer',
        '&:hover': {
          backgroundColor: 'primary.light',
          color: 'white',
        },
      }}
    />
  );
};
