// 管理画面共通JavaScript

// クリップボードにコピー
function copyToClipboard(text, buttonElement) {
    if (typeof text === 'string') {
        // 文字列が直接渡された場合
        navigator.clipboard.writeText(text).then(() => {
            if (buttonElement) {
                const originalText = buttonElement.textContent;
                buttonElement.textContent = 'コピーしました！';
                buttonElement.classList.add('is-success');

                setTimeout(() => {
                    buttonElement.textContent = originalText;
                    buttonElement.classList.remove('is-success');
                }, 2000);
            } else {
                alert('コピーしました！');
            }
        }).catch(err => {
            alert('コピーに失敗しました');
        });
    } else {
        // 要素IDが渡された場合
        const element = document.getElementById(text);
        if (element) {
            const value = element.value || element.textContent;
            navigator.clipboard.writeText(value).then(() => {
                alert('コピーしました！');
            }).catch(err => {
                alert('コピーに失敗しました');
            });
        }
    }
}

// 確認ダイアログ
function confirmAction(message, callback) {
    if (confirm(message)) {
        callback();
    }
}

// 日付フォーマット
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('ja-JP', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

// 数値フォーマット
function formatNumber(num) {
    return num.toLocaleString('ja-JP');
}
