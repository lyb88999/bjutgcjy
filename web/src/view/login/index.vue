<template>
  <div id="userLayout" class="w-full h-full relative">
    <div
      class="rounded-lg flex items-center justify-evenly w-full h-full bg-white md:w-screen md:h-screen md:bg-[#194bfb]">
      <div class="md:bg-[rgb(161,191,244)] w-3/5 w-10/12 h-full flex items-center justify-evenly">
        <div class="oblique h-[130%] w-3/5 bg-white transform -rotate-12 absolute -ml-52" />
        <!-- 分割斜块 -->
        <div class="z-[999] pt-12 pb-10 md:w-96 w-full  rounded-lg flex flex-col justify-between box-border">
          <div>
            <div class="flex items-center justify-center">

              <img class="w-24" src="@/assets/logoNew.png" alt>
            </div>
            <div class="mb-9">
              <p class="text-center text-3xl font-bold">{{ $GIN_VUE_ADMIN.appName }}</p>
              <p class="text-center text-sm font-normal text-gray-500 mt-2.5">北京工业大学采集点
              </p>
            </div>
            <el-form ref="loginForm" :model="loginFormData" :rules="rules" :validate-on-rule-change="false" @keyup.enter="submitForm">
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
                  <!-- <template #suffix>
                  </template> -->
                </el-input>
              </el-form-item>
              <el-form-item prop="referrer_username" v-if="loginType.status">
                <el-input v-model="loginFormData.referrer_username"
                  placeholder="请输入邀请人">
                  <!-- <template #suffix>
                  </template> -->
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
  // return await userRegister(loginFormData)
  return await userStore.Register(loginFormData)
}
const submitForm = () => {
  loginForm.value.validate(async (v) => {
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
