function registerAuthComponent() {
    if (window.AuthComponentRegistered) return;
    window.AuthComponentRegistered = true;

    Alpine.data('authHandler', () => ({
        mainMode: 'auth',
        activeTab: 'registration',
        activeApiTab: 'list',
        isJwtAuth: false,
        isSessionAuth: false,
        isApiKeyAuth: false,
        isBasicAuth: false,
        isOAuthAuth: false,
        
        async checkAuth(url) {
            try {
                const response = await window.apiFetch(url, { 
                    method: 'POST', 
                    credentials: 'include',
                    headers: window.getAuthHeaders()
                });
                return response.ok;
            } catch (e) {
                return false;
            }
        },
        
        async updateAuthStatus() {
            switch(this.activeTab) {
                case 'jwt':
                    this.isJwtAuth = await this.checkAuth('/auth/jwt/validate');
                    break;
                case 'oauth':
                    this.isOAuthAuth = await this.checkAuth('/auth/oauth/validate');
                    break;
                case 'basic':
                    this.isBasicAuth = await this.checkAuth('/auth/basic/validate');
                    break;
                case 'session':
                    this.isSessionAuth = await this.checkAuth('/auth/session/validate');
                    break;
                case 'apikey':
                    this.isApiKeyAuth = await this.checkAuth('/auth/apikey/validate');
                    break;
            }
        },
        
        async logout() {
            let endpoint = '';
            switch(this.activeTab) {
                case 'jwt':
                    endpoint = '/auth/jwt/logout';
                    break;
                case 'oauth':
                    endpoint = '/auth/oauth/logout';
                    break;
                case 'basic':
                    const u = document.getElementById('basic-username');
                    const p = document.getElementById('basic-password');
                    if(u) u.value = '';
                    if(p) p.value = '';
                    break;
                case 'session':
                    endpoint = '/auth/session/revoke';
                    break;
                case 'apikey':
                    endpoint = '/auth/apikey/revoke';
                    break;
            }

            if (endpoint) {
                try {
                    await fetch(endpoint, { method: 'POST', credentials: 'include', headers: window.getAuthHeaders() });
                } catch (e) {}
            }

            sessionStorage.removeItem('basic_token');
            await this.updateAuthStatus();
            this.$dispatch('notify', { message: 'Выход выполнен успешно', type: 'info' });
        },
        
        init() {
            this.updateAuthStatus();
            const clearResults = () => {
                const authRes = document.getElementById('auth-result');
                if (authRes) authRes.innerHTML = '';
            };
            this.$watch('mainMode', clearResults);
            this.$watch('activeApiTab', clearResults);
            this.$watch('activeTab', () => {
                this.updateAuthStatus();
                clearResults();
            });
            window.addEventListener('auth-updated', () => this.updateAuthStatus());
        }
    }));
}

document.addEventListener('alpine:init', registerAuthComponent);
if (window.Alpine) registerAuthComponent();
