import { useMutation, useQueryClient } from '@tanstack/react-query';
import { uploadMedia, deleteMedia } from '../api/media';

/**
 * メディアアップロード
 */
export const useUploadMedia = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ postId, files }: { postId: number; files: File[] }) =>
      uploadMedia(postId, files),
    onSuccess: () => {
      // 投稿一覧を再取得（メディアが更新されたため）
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['post'] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
    },
  });
};

/**
 * メディア削除
 */
export const useDeleteMedia = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteMedia,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['post'] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
    },
  });
};
