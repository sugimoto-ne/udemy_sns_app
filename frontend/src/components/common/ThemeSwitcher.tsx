import { useState } from 'react';
import {
  IconButton,
  Menu,
  MenuItem,
  ListItemIcon,
  ListItemText,
  Tooltip,
} from '@mui/material';
import PaletteIcon from '@mui/icons-material/Palette';
import CheckIcon from '@mui/icons-material/Check';
import { useTheme } from '../../contexts/ThemeContext';
import { themeDisplayNames } from '../../theme/themes';
import type { ThemeName } from '../../theme/themes';

export const ThemeSwitcher = () => {
  const { currentTheme, setTheme } = useTheme();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);

  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleThemeSelect = (theme: ThemeName) => {
    setTheme(theme);
    handleClose();
  };

  const themeOptions: ThemeName[] = ['default', 'dark', 'purple', 'green', 'orange'];

  return (
    <>
      <Tooltip title="テーマを変更">
        <IconButton
          onClick={handleClick}
          size="large"
          aria-controls={open ? 'theme-menu' : undefined}
          aria-haspopup="true"
          aria-expanded={open ? 'true' : undefined}
          color="inherit"
        >
          <PaletteIcon />
        </IconButton>
      </Tooltip>
      <Menu
        id="theme-menu"
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
        MenuListProps={{
          'aria-labelledby': 'theme-button',
        }}
      >
        {themeOptions.map((theme) => (
          <MenuItem
            key={theme}
            onClick={() => handleThemeSelect(theme)}
            selected={currentTheme === theme}
          >
            <ListItemIcon>
              {currentTheme === theme && <CheckIcon fontSize="small" />}
            </ListItemIcon>
            <ListItemText>{themeDisplayNames[theme]}</ListItemText>
          </MenuItem>
        ))}
      </Menu>
    </>
  );
};
