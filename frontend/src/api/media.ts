import { apiClient } from './client';

export interface Media {
  id: number;
  post_id: number;
  media_type: 'image' | 'video' | 'audio';
  media_url: string;
  file_size: number;
  duration?: number;
  order_index: number;
  created_at: string;
}

export interface MediaUploadResponse {
  message: string;
  media: Media[];
}

/**
 * メディアをアップロード
 */
export const uploadMedia = async (
  postId: number,
  files: File[]
): Promise<MediaUploadResponse> => {
  const formData = new FormData();
  formData.append('post_id', postId.toString());

  files.forEach((file) => {
    formData.append('files', file);
  });

  const response = await apiClient.post('/media/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });

  return response.data;
};

/**
 * メディアを削除
 */
export const deleteMedia = async (mediaId: number): Promise<{ message: string }> => {
  const response = await apiClient.delete(`/media/${mediaId}`);
  return response.data;
};
