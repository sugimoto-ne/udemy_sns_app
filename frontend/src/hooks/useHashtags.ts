import { useQuery, useInfiniteQuery } from '@tanstack/react-query';
import { getTrendingHashtags, getPostsByHashtag } from '../api/hashtags';

/**
 * トレンドハッシュタグ取得
 */
export const useTrendingHashtags = (limit: number = 10) => {
  return useQuery({
    queryKey: ['hashtags', 'trending', limit],
    queryFn: () => getTrendingHashtags({ limit }),
  });
};

/**
 * ハッシュタグ別投稿取得（無限スクロール対応）
 */
export const useHashtagPosts = (hashtagName: string) => {
  return useInfiniteQuery({
    queryKey: ['hashtags', hashtagName, 'posts'],
    queryFn: ({ pageParam }) =>
      getPostsByHashtag(hashtagName, {
        limit: 20,
        cursor: pageParam,
      }),
    getNextPageParam: (lastPage) => {
      return lastPage.pagination.has_more ? lastPage.pagination.next_cursor : undefined;
    },
    initialPageParam: undefined as string | undefined,
    enabled: !!hashtagName,
  });
};
