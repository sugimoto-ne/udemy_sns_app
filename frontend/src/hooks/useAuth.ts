import { useMutation } from '@tanstack/react-query';
import { requestPasswordReset, confirmPasswordReset } from '../api/password-reset';
import { verifyEmail, resendVerificationEmail } from '../api/email-verification';

/**
 * パスワードリセットリクエスト
 */
export const useRequestPasswordReset = () => {
  return useMutation({
    mutationFn: requestPasswordReset,
  });
};

/**
 * パスワードリセット確認
 */
export const useConfirmPasswordReset = () => {
  return useMutation({
    mutationFn: confirmPasswordReset,
  });
};

/**
 * メールアドレス認証
 */
export const useVerifyEmail = () => {
  return useMutation({
    mutationFn: verifyEmail,
  });
};

/**
 * 認証メール再送信
 */
export const useResendVerificationEmail = () => {
  return useMutation({
    mutationFn: resendVerificationEmail,
  });
};
