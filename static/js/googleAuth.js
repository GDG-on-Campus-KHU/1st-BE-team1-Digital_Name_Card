// googleAuth.js

class AuthService {
  constructor() {
    this.token = null;
    this.init();
  }

  init() {
    this.token = this.getCookie('token');
    if (this.token) {
      console.log('토큰 발견:', this.token);
      this.setupAuthenticatedState();
    }
  }

  getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
  }

  setupAuthenticatedState() {
    // Authorization 헤더에 JWT 토큰을 기본으로 포함
    this.setupAxiosInterceptors();
    this.updateUIForAuthenticatedUser();
  }

  setupAxiosInterceptors() {
    axios.interceptors.request.use(
      (config) => {
        if (this.token) {
          config.headers.Authorization = `Bearer ${this.token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // 토큰 만료 등의 인증 오류 처리
    axios.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response.status === 401) {
          this.handleAuthError();
        }
        return Promise.reject(error);
      }
    );
  }

  async handleGoogleLogin() {
    try {
      window.location.href = '/auth/google/login';
    } catch (error) {
      console.error('Google 로그인 시도 실패:', error);
      this.showError('Google 로그인 연결에 실패했습니다. 잠시 후 다시 시도해주세요.');
    }
  }

  async fetchUserProfile() {
    try {
      const response = await axios.get('/profile');
      return response.data;
    } catch (error) {
      console.error('프로필 정보 가져오기 실패:', error);
      throw error;
    }
  }

  updateUIForAuthenticatedUser() {
    this.fetchUserProfile()
      .then(userData => {
        const userEmail = document.getElementById('userEmail');
        const userName = document.getElementById('userName');
        
        if (userEmail) userEmail.textContent = userData.email || '';
        if (userName) userName.textContent = userData.nickname || '';
        
        // 로그인 버튼을 로그아웃 버튼으로 변경
        this.updateLoginButton();
      })
      .catch(error => {
        this.handleAuthError();
      });
  }

  updateLoginButton() {
    const loginButton = document.getElementById('google-login-button');
    if (loginButton) {
      loginButton.textContent = '로그아웃';
      loginButton.onclick = () => this.handleLogout();
    }
  }

  handleAuthError() {
    // 인증 에러 발생 시 (토큰 만료 등) 쿠키 삭제 및 재로그인 필요
    document.cookie = 'token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    this.token = null;
    window.location.href = '/';
  }

  handleLogout() {
    document.cookie = 'token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    this.token = null;
    window.location.reload();
  }

  showError(message) {
    const resultDiv = document.getElementById('google-result');
    if (resultDiv) {
      resultDiv.innerHTML = `<p class="error">${message}</p>`;
    }
  }
}

// 초기화 및 이벤트 리스너 설정
document.addEventListener('DOMContentLoaded', () => {
  const authService = new AuthService();
  
  const googleLoginButton = document.getElementById('google-login-button');
  if (googleLoginButton) {
    googleLoginButton.addEventListener('click', () => {
      if (authService.token) {
        authService.handleLogout();
      } else {
        authService.handleGoogleLogin();
      }
    });
  }
});