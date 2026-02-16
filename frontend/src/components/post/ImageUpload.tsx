import React, { useState, useEffect } from 'react';
import {
  Box,
  IconButton,
  ImageList,
  ImageListItem,
  ImageListItemBar,
  Button,
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import AddPhotoAlternateIcon from '@mui/icons-material/AddPhotoAlternate';

interface ImageUploadProps {
  maxFiles?: number;
  onFilesChange: (files: File[]) => void;
  selectedFiles?: File[]; // controlled component用
}

export const ImageUpload: React.FC<ImageUploadProps> = ({
  maxFiles = 4,
  onFilesChange,
  selectedFiles: externalSelectedFiles = []
}) => {
  const [selectedFiles, setSelectedFiles] = useState<File[]>(externalSelectedFiles);
  const [previewUrls, setPreviewUrls] = useState<string[]>([]);

  // 外部から渡されたselectedFilesが変更されたら同期
  useEffect(() => {
    // 外部からリセットされた場合（空配列になった場合）
    if (externalSelectedFiles.length === 0 && selectedFiles.length > 0) {
      // 既存のプレビューURLを全て解放
      previewUrls.forEach(url => URL.revokeObjectURL(url));
      setPreviewUrls([]);
      setSelectedFiles([]);
    }
  }, [externalSelectedFiles]);

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(event.target.files || []);
    const remainingSlots = maxFiles - selectedFiles.length;

    if (files.length > remainingSlots) {
      alert(`最大${maxFiles}枚までアップロードできます`);
      return;
    }

    const newFiles = [...selectedFiles, ...files].slice(0, maxFiles);
    setSelectedFiles(newFiles);
    onFilesChange(newFiles);

    // プレビュー画像を生成
    const newUrls = files.map((file) => URL.createObjectURL(file));
    setPreviewUrls((prev) => [...prev, ...newUrls]);
  };

  const handleRemove = (index: number) => {
    const newFiles = selectedFiles.filter((_, i) => i !== index);
    const newUrls = previewUrls.filter((_, i) => i !== index);

    // 古いURLを解放
    URL.revokeObjectURL(previewUrls[index]);

    setSelectedFiles(newFiles);
    setPreviewUrls(newUrls);
    onFilesChange(newFiles);
  };

  return (
    <Box>
      {previewUrls.length > 0 && (
        <ImageList cols={Math.min(previewUrls.length, 2)} gap={8} sx={{ mb: 2 }}>
          {previewUrls.map((url, index) => (
            <ImageListItem key={index}>
              <img
                src={url}
                alt={`Preview ${index + 1}`}
                style={{ height: 200, objectFit: 'cover' }}
              />
              <ImageListItemBar
                actionIcon={
                  <IconButton
                    sx={{ color: 'white' }}
                    onClick={() => handleRemove(index)}
                    size="small"
                  >
                    <DeleteIcon />
                  </IconButton>
                }
              />
            </ImageListItem>
          ))}
        </ImageList>
      )}

      {selectedFiles.length < maxFiles && (
        <Button
          component="label"
          variant="outlined"
          startIcon={<AddPhotoAlternateIcon />}
          fullWidth
        >
          画像を追加 ({selectedFiles.length}/{maxFiles})
          <input
            type="file"
            hidden
            multiple
            accept="image/*"
            onChange={handleFileSelect}
          />
        </Button>
      )}
    </Box>
  );
};
