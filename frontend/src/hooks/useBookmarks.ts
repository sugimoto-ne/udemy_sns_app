import { useMutation, useQueryClient, useInfiniteQuery } from '@tanstack/react-query';
import { bookmarkPost, unbookmarkPost, getBookmarks } from '../api/bookmarks';

/**
 * ブックマーク追加
 */
export const useBookmarkPost = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: bookmarkPost,
    onMutate: async (postId) => {
      // 進行中のクエリをキャンセル
      await queryClient.cancelQueries({ queryKey: ['timeline'] });

      // 現在のデータを取得（ロールバック用）
      const previousData = queryClient.getQueryData(['timeline']);

      // 楽観的更新 - タイムラインのis_bookmarkedをtrueに
      queryClient.setQueriesData({ queryKey: ['timeline'] }, (old: any) => {
        if (!old?.data) return old;
        return {
          ...old,
          data: old.data.map((post: any) =>
            post.id === postId ? { ...post, is_bookmarked: true } : post
          ),
        };
      });

      return { previousData };
    },
    onError: (_err, _postId, context) => {
      // エラー時はロールバック
      if (context?.previousData) {
        queryClient.setQueryData(['timeline'], context.previousData);
      }
    },
    onSettled: () => {
      // 完了後に再取得して同期
      queryClient.invalidateQueries({ queryKey: ['bookmarks'] });
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['post'] });
    },
  });
};

/**
 * ブックマーク解除
 */
export const useUnbookmarkPost = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: unbookmarkPost,
    onMutate: async (postId) => {
      await queryClient.cancelQueries({ queryKey: ['timeline'] });
      const previousData = queryClient.getQueryData(['timeline']);

      // 楽観的更新 - タイムラインのis_bookmarkedをfalseに
      queryClient.setQueriesData({ queryKey: ['timeline'] }, (old: any) => {
        if (!old?.data) return old;
        return {
          ...old,
          data: old.data.map((post: any) =>
            post.id === postId ? { ...post, is_bookmarked: false } : post
          ),
        };
      });

      return { previousData };
    },
    onError: (_err, _postId, context) => {
      if (context?.previousData) {
        queryClient.setQueryData(['timeline'], context.previousData);
      }
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['bookmarks'] });
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['post'] });
    },
  });
};

/**
 * ブックマーク一覧取得（無限スクロール対応）
 */
export const useBookmarks = () => {
  return useInfiniteQuery({
    queryKey: ['bookmarks'],
    queryFn: ({ pageParam }) =>
      getBookmarks({
        limit: 20,
        cursor: pageParam,
      }),
    getNextPageParam: (lastPage) => {
      return lastPage.pagination.has_more ? lastPage.pagination.next_cursor : undefined;
    },
    initialPageParam: undefined as string | undefined,
  });
};
