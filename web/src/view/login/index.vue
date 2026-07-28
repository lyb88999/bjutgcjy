<template>
  <div id="userLayout" class="w-full h-full relative">
    <!-- 品牌切换按钮：两套入口共用同一套账号体系和登录逻辑，只是登录页视觉不同，
         方便教育数据库和专家库两边的用户都能在自己熟悉的登录页登录 -->
    <button type="button" class="brand-switch" @click="toggleBrand">
      切换到{{ brandTheme === 'expert' ? '教育数据库' : '专家库' }}入口
    </button>

    <!-- 专家库入口：参考北京市社科联官网（bjsk.org.cn）色调重做的藏青+中国红风格 -->
    <div v-if="brandTheme === 'expert'" id="expertLayout" class="login-page">
      <div class="login-left">
        <div class="login-form-wrap">
          <div class="brand">
            <div class="brand-mark">
              <span class="brand-mark-bar bar-1" />
              <span class="brand-mark-bar bar-2" />
              <span class="brand-mark-bar bar-3" />
            </div>
            <div class="brand-text">
              <p class="brand-title">{{ $GIN_VUE_ADMIN.appName }}</p>
              <p class="brand-subtitle">智库专家资源管理与推荐平台</p>
            </div>
          </div>

          <el-form ref="loginForm" class="login-form" :model="loginFormData" :rules="rules" :validate-on-rule-change="false" @keyup.enter="submitForm">
            <el-form-item prop="username">
              <el-input v-model="loginFormData.username" placeholder="请输入用户名">
                <template #suffix>
                  <span class="input-icon">
                    <el-icon>
                      <user />
                    </el-icon>
                  </span>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item prop="password">
              <el-input v-model="loginFormData.password" :type="lock === 'lock' ? 'password' : 'text'"
                placeholder="请输入密码">
                <template #suffix>
                  <span class="input-icon">
                    <el-icon>
                      <component :is="lock" @click="changeLock" />
                    </el-icon>
                  </span>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item prop="invite_code" v-if="loginType.status">
              <el-input v-model="loginFormData.invite_code"
                placeholder="请输入邀请码">
              </el-input>
            </el-form-item>
            <el-form-item prop="referrer_username" v-if="loginType.status">
              <el-input v-model="loginFormData.referrer_username"
                placeholder="请输入邀请人">
              </el-input>
            </el-form-item>
            <el-form-item prop="captcha">
              <div class="vPicBox">
                <el-input v-model="loginFormData.captcha" placeholder="请输入验证码" style="width: 60%" />
                <div class="vPic">
                  <img v-if="picPath" :src="picPath" alt="请输入验证码" @click="loginVerify()">
                </div>
              </div>
            </el-form-item>
            <el-form-item>
              <el-button class="btn-ghost" style="width: 38%" size="large" @click="checkInit">前往初始化</el-button>
              <el-button class="btn-primary" size="large" style="width: 38%; margin-left: 8%" @click="submitForm">
                <div v-if="loginType.status">注 册</div>
                <div v-if="!loginType.status">登 录</div>
              </el-button>
              <el-switch v-model="loginType.status" style="width: 13%; margin-left: 3%" />
            </el-form-item>
          </el-form>
        </div>
      </div>

      <div class="login-right hidden md:block">
        <svg class="rings" viewBox="0 0 600 600" xmlns="http://www.w3.org/2000/svg">
          <circle cx="420" cy="300" r="260" />
          <circle cx="420" cy="300" r="200" />
          <circle cx="420" cy="300" r="140" />
        </svg>
        <svg class="lattice" viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
          <pattern id="latticePattern" width="40" height="40" patternUnits="userSpaceOnUse">
            <path d="M0 20 H40 M20 0 V40" stroke="rgba(255,255,255,0.10)" stroke-width="1" />
            <path d="M0 0 L40 40 M40 0 L0 40" stroke="rgba(255,255,255,0.06)" stroke-width="1" />
          </pattern>
          <rect width="200" height="200" fill="url(#latticePattern)" />
        </svg>

        <div class="network-graphic">
          <svg viewBox="0 0 400 300" xmlns="http://www.w3.org/2000/svg">
            <g stroke="rgba(255,255,255,0.35)" stroke-width="1.2">
              <line x1="80" y1="90" x2="200" y2="60" />
              <line x1="80" y1="90" x2="150" y2="180" />
              <line x1="200" y1="60" x2="320" y2="100" />
              <line x1="200" y1="60" x2="150" y2="180" />
              <line x1="150" y1="180" x2="260" y2="220" />
              <line x1="320" y1="100" x2="260" y2="220" />
              <line x1="150" y1="180" x2="80" y2="230" />
            </g>
            <g fill="#fff">
              <circle cx="80" cy="90" r="6" />
              <circle cx="200" cy="60" r="8" />
              <circle cx="320" cy="100" r="6" />
              <circle cx="150" cy="180" r="9" />
              <circle cx="260" cy="220" r="6" />
              <circle cx="80" cy="230" r="5" />
            </g>
            <g fill="#d4a017">
              <circle cx="150" cy="180" r="4" />
            </g>
          </svg>
        </div>

        <div class="right-caption">
          <div class="right-ribbon">汇智专家资源 · 支撑首都决策</div>
          <p class="right-desc">背景信息 · 学科画像 · 研究成果 · 决策影响 · 综合排序推荐</p>
        </div>
      </div>
    </div>

    <!-- 教育数据库入口：保留原有的蓝色斜切风格与"北京工业大学采集点"品牌文案 -->
    <div v-else class="rounded-lg flex items-center justify-evenly w-full h-full bg-white md:w-screen md:h-screen md:bg-[#194bfb]">
      <div class="md:bg-[rgb(161,191,244)] w-3/5 w-10/12 h-full flex items-center justify-evenly">
        <div class="oblique h-[130%] w-3/5 bg-white transform -rotate-12 absolute -ml-52" />
        <!-- 分割斜块 -->
        <div class="z-[999] pt-12 pb-10 md:w-96 w-full  rounded-lg flex flex-col justify-between box-border">
          <div>
            <div class="flex items-center justify-center">
              <img class="w-24" src="@/assets/logoNew.png" alt>
            </div>
            <div class="mb-9">
              <p class="text-center text-3xl font-bold">工程教育数据库</p>
              <p class="text-center text-sm font-normal text-gray-500 mt-2.5">北京工业大学采集点
              </p>
            </div>
            <el-form ref="loginFormEdu" :model="loginFormData" :rules="rules" :validate-on-rule-change="false" @keyup.enter="submitForm">
              <el-form-item prop="username">
                <el-input v-model="loginFormData.username" placeholder="请输入用户名">
                  <template #suffix>
                    <span class="input-icon">
                      <el-icon>
                        <user />
                      </el-icon>
                    </span>
                  </template>
                </el-input>
              </el-form-item>
              <el-form-item prop="password">
                <el-input v-model="loginFormData.password" :type="lock === 'lock' ? 'password' : 'text'"
                  placeholder="请输入密码">
                  <template #suffix>
                    <span class="input-icon">
                      <el-icon>
                        <component :is="lock" @click="changeLock" />
                      </el-icon>
                    </span>
                  </template>
                </el-input>
              </el-form-item>
              <el-form-item prop="invite_code" v-if="loginType.status">
                <el-input v-model="loginFormData.invite_code"
                  placeholder="请输入邀请码">
                </el-input>
              </el-form-item>
              <el-form-item prop="referrer_username" v-if="loginType.status">
                <el-input v-model="loginFormData.referrer_username"
                  placeholder="请输入邀请人">
                </el-input>
              </el-form-item>
              <el-form-item prop="captcha">
                <div class="vPicBox">
                  <el-input v-model="loginFormData.captcha" placeholder="请输入验证码" style="width: 60%" />
                  <div class="vPic">
                    <img v-if="picPath" :src="picPath" alt="请输入验证码" @click="loginVerify()">
                  </div>
                </div>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" style="width: 38%" size="large" @click="checkInit">前往初始化</el-button>
                <el-button type="primary" size="large" style="width: 38%; margin-left: 8%" @click="submitForm">
                  <div v-if="loginType.status">注 册</div>
                  <div v-if="!loginType.status">登 录</div>
                </el-button>
                <el-switch v-model="loginType.status" style="width: 13%; margin-left: 3%" />
              </el-form-item>
            </el-form>
          </div>
        </div>
      </div>
      <div class="hidden md:block w-1/2 h-full float-right bg-[#194bfb]"><img class="h-full" src="@/assets/banner.jpg"
          alt="banner"></div>
    </div>

    <BottomInfo class="left-0 right-0 absolute bottom-3 mx-auto  w-full z-20">
      <div class="links items-center justify-center gap-2 hidden md:flex">
        <a href="http://doc.henrongyi.top/" target="_blank">
          <img src="@/assets/docs.png" class="w-8 h-8" alt="文档">
        </a>
        <a href="https://support.qq.com/product/371961" target="_blank">
          <img src="@/assets/kefu.png" class="w-8 h-8" alt="客服">
        </a>
        <a href="https://github.com/flipped-aurora/gin-vue-admin" target="_blank">
          <img src="@/assets/github.png" class="w-8 h-8" alt="github">
        </a>
        <a href="https://space.bilibili.com/322210472" target="_blank">
          <img src="@/assets/video.png" class="w-8 h-8" alt="视频站">
        </a>
      </div>
    </BottomInfo>
  </div>
</template>

<script setup>
import { captcha } from '@/api/user'
import { checkDB } from '@/api/initdb'
import BottomInfo from '@/view/layout/bottomInfo/bottomInfo.vue'
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({
  name: 'Login',
})

const router = useRouter()

// 品牌入口切换：教育数据库和专家库共用同一套登录逻辑/账号体系，只是登录页外观不同，
// 记住上次选择，方便各自的用户回来时不用重新切换
const BRAND_STORAGE_KEY = 'gva-login-brand-theme'
const brandTheme = ref(localStorage.getItem(BRAND_STORAGE_KEY) || 'expert')
const toggleBrand = () => {
  brandTheme.value = brandTheme.value === 'expert' ? 'edu' : 'expert'
  localStorage.setItem(BRAND_STORAGE_KEY, brandTheme.value)
}

// 验证函数
const checkUsername = (rule, value, callback) => {
  if (value.length < 5) {
    return callback(new Error('请输入正确的用户名'))
  } else {
    callback()
  }
}
const checkPassword = (rule, value, callback) => {
  if (value.length < 6) {
    return callback(new Error('请输入正确的密码'))
  } else {
    callback()
  }
}

// 获取验证码
const loginVerify = () => {
  captcha({}).then(async (ele) => {
    rules.captcha.push({
      max: ele.data.captchaLength,
      min: ele.data.captchaLength,
      message: `请输入${ele.data.captchaLength}位验证码`,
      trigger: 'blur',
    })
    picPath.value = ele.data.picPath
    loginFormData.captchaId = ele.data.captchaId
    loginFormData.openCaptcha = ele.data.openCaptcha
  })
}
loginVerify()

// 登录相关操作
const lock = ref('lock')
const changeLock = () => {
  lock.value = lock.value === 'lock' ? 'unlock' : 'lock'
}

const loginForm = ref(null)
const loginFormEdu = ref(null)
const picPath = ref('')
const loginFormData = reactive({
  username: '',
  password: '',
  captcha: '',
  captchaId: '',
  invite_code: '',
  referrer_username: ''
})
const loginType = reactive({
  status: false,
})
const rules = reactive({
  username: [{ validator: checkUsername, trigger: 'blur' }],
  password: [{ validator: checkPassword, trigger: 'blur' }],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    {
      message: '验证码格式不正确',
      trigger: 'blur',
    },
  ],
})

const userStore = useUserStore()
const login = async () => {
  return await userStore.LoginIn(loginFormData)
}
const register = async () => {
  return await userStore.Register(loginFormData)
}
const submitForm = () => {
  const activeForm = brandTheme.value === 'expert' ? loginForm.value : loginFormEdu.value
  activeForm.validate(async (v) => {
    if (v) {
      let flag
      if (loginType.status) {
        flag = await register()
      } else {
        flag = await login()
      }
      if (!flag) {
        loginVerify()
        loginType.status = false
      }
    } else {
      ElMessage({
        type: 'error',
        message: '请正确填写登录信息',
        showClose: true,
      })
      loginVerify()
      return false
    }
  })
}

// 跳转初始化
const checkInit = async () => {
  const res = await checkDB()
  if (res.code === 0) {
    if (res.data?.needInit) {
      userStore.NeedInit()
      router.push({ name: 'Init' })
    } else {
      ElMessage({
        type: 'info',
        message: '已配置数据库信息，无法初始化',
      })
    }
  }
}

</script>

<style lang="scss" scoped>
// 配色取自北京市社科联官网（bjsk.org.cn）的视觉基调：深藏青蓝 + 中国红 + 淡青蓝水纹，
// 传达"官方智库/社科系统"的严肃与专业感，区别于 gin-vue-admin 默认的通用后台配色
$navy: #16407a;
$navy-dark: #0c2a54;
$red: #b5272d;
$gold: #d4a017;
$ice: #eaf4fb;

.brand-switch {
  position: absolute;
  top: 16px;
  right: 20px;
  z-index: 1000;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 999px;
  padding: 6px 16px;
  font-size: 12px;
  color: #333;
  cursor: pointer;

  &:hover {
    background: #fff;
  }
}

#expertLayout.login-page {
  width: 100%;
  height: 100%;
  display: flex;
  position: relative;
  border-top: 4px solid $red;
  box-sizing: border-box;

  @media (max-width: 767px) {
    border-top: none;
  }
}

.login-left {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(180deg, $ice 0%, #ffffff 40%);

  @media (min-width: 768px) {
    width: 44%;
  }
}

.login-form-wrap {
  width: 100%;
  max-width: 380px;
  padding: 0 24px;
  box-sizing: border-box;
}

.brand {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 40px;
}

.brand-mark {
  width: 52px;
  height: 52px;
  flex: none;
  border-radius: 10px;
  background: linear-gradient(135deg, $red 0%, #8f1f22 100%);
  box-shadow: 0 4px 10px rgba(181, 39, 45, 0.35);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.brand-mark-bar {
  display: block;
  height: 3px;
  border-radius: 2px;
  background: #fff;
}

.brand-mark .bar-1 { width: 26px; }
.brand-mark .bar-2 { width: 20px; }
.brand-mark .bar-3 { width: 14px; background: $gold; }

.brand-title {
  font-size: 21px;
  font-weight: 700;
  color: $navy-dark;
  line-height: 1.3;
  margin: 0;
}

.brand-subtitle {
  font-size: 12px;
  color: #7a8aa0;
  margin: 4px 0 0;
}

.login-form {
  :deep(.el-input__wrapper) {
    box-shadow: 0 0 0 1px #d7e2ee inset;
  }
  :deep(.el-input__wrapper.is-focus) {
    box-shadow: 0 0 0 1px $navy inset;
  }
}

.btn-primary {
  background: $navy;
  border-color: $navy;
  color: #fff;

  &:hover, &:focus {
    background: lighten($navy, 8%);
    border-color: lighten($navy, 8%);
  }
}

.btn-ghost {
  background: #fff;
  border-color: $navy;
  color: $navy;

  &:hover, &:focus {
    background: $ice;
    border-color: $navy;
    color: $navy;
  }
}

.vPicBox {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;

  .vPic {
    width: 38%;
    height: 40px;
    cursor: pointer;

    img {
      width: 100%;
      height: 100%;
      border-radius: 4px;
    }
  }
}

.login-right {
  position: relative;
  width: 56%;
  height: 100%;
  overflow: hidden;
  background: linear-gradient(150deg, $navy 0%, $navy-dark 100%);
  display: flex;
  align-items: center;
  justify-content: center;

  .rings {
    position: absolute;
    right: -120px;
    top: 50%;
    transform: translateY(-50%);
    width: 640px;
    height: 640px;
    fill: none;
    stroke: rgba(255, 255, 255, 0.14);
    stroke-width: 1.5;
  }

  .lattice {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
  }

  .network-graphic {
    position: relative;
    width: 420px;
    max-width: 70%;
  }

  .right-caption {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 64px;
    text-align: center;
    padding: 0 40px;
  }

  .right-ribbon {
    display: inline-block;
    background: $red;
    color: #fff;
    font-size: 15px;
    font-weight: 600;
    letter-spacing: 1px;
    padding: 8px 22px;
    border-radius: 2px;
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.25);
  }

  .right-desc {
    margin-top: 14px;
    color: rgba(255, 255, 255, 0.75);
    font-size: 13px;
    letter-spacing: 0.5px;
  }
}
</style>
