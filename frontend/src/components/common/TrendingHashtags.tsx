import React from 'react';
import {
  Box,
  Typography,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Paper,
  CircularProgress,
} from '@mui/material';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import { useTrendingHashtags } from '../../hooks/useHashtags';
import { useNavigate } from 'react-router-dom';

export const TrendingHashtags: React.FC = () => {
  const { data, isLoading, error } = useTrendingHashtags(10);
  const navigate = useNavigate();

  if (isLoading) {
    return (
      <Paper sx={{ p: 2 }}>
        <Box display="flex" justifyContent="center">
          <CircularProgress size={24} />
        </Box>
      </Paper>
    );
  }

  if (error || !data?.hashtags || data.hashtags.length === 0) {
    return null;
  }

  return (
    <Paper sx={{ p: 2 }}>
      <Box display="flex" alignItems="center" gap={1} mb={2}>
        <TrendingUpIcon color="primary" />
        <Typography variant="h6">トレンド</Typography>
      </Box>

      <List disablePadding>
        {data.hashtags.map((hashtag) => (
          <ListItem key={hashtag.id} disablePadding>
            <ListItemButton
              onClick={() => navigate(`/hashtags/${encodeURIComponent(hashtag.name)}`)}
            >
              <ListItemText
                primary={`#${hashtag.name}`}
                secondary={`${hashtag.count}件の投稿`}
                primaryTypographyProps={{
                  fontWeight: 'bold',
                  color: 'primary',
                }}
              />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Paper>
  );
};
