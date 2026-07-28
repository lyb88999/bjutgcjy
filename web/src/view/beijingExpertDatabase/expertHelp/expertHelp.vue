<template>
  <div class="expert-help">
    <div class="hero">
      <h2>北京市哲学社会科学专家库</h2>
      <p>这是一个记录、审核、查找全市哲学社会科学领域专家的系统。专家信息填好之后，要经过审核才会正式收录；收录之后，别人就能通过关键词搜到这位专家、了解他的研究方向和过往成果，方便找到合适的人参与决策咨询。</p>
    </div>

    <div class="role-grid">
      <div v-for="role in roleCards" :key="role.key" class="role-card" :class="{ active: role.key === myRoleKey }">
        <el-icon class="role-icon"><component :is="role.icon" /></el-icon>
        <div class="role-name">
          {{ role.name }}
          <el-tag v-if="role.key === myRoleKey" type="success" size="small" style="margin-left: 6px;">这是你</el-tag>
        </div>
        <div class="role-desc">{{ role.desc }}</div>
      </div>
    </div>

    <el-tabs v-model="activeTab" class="help-tabs">
      <el-tab-pane label="怎么登录" name="login">
        <p>打开系统，输入账号、密码，再输入图形验证码，点"登录"就行。</p>
        <p>如果登录进去发现左边菜单是空的、或者少了一些功能，先别慌，大概率是账号的角色或者所在单位还没配置好，找负责管理系统的人确认一下。</p>
        <p style="color:#909399;">如果刚被调整过权限，光刷新页面有时候不会立刻生效，退出重新登录一次就好了。</p>
      </el-tab-pane>

      <el-tab-pane v-if="isApplicant || isAdmin" label="填报专家信息" name="applicant">
        <el-steps :active="4" align-center finish-status="success" style="margin-bottom: 24px;">
          <el-step title="填写背景信息" description="姓名、单位、职称等基本资料" />
          <el-step title="补充研究成果" description="成果、决策影响、学术兼职、标签" />
          <el-step title="提交审核" description="一键提交，等待审核结果" />
          <el-step title="正式收录" description="审核通过后可以被搜索到" />
        </el-steps>

        <h4>第一步：新增专家档案</h4>
        <p>在"专家主档"页面点"新增"，把姓名、所在单位、职称、学历、研究方向这些基本信息填好。</p>
        <p><strong>所在单位</strong>是一个可以搜索的下拉框，优先从列表里选已经有的单位；如果确实找不到你的单位，直接输入名字保存也可以，之后系统管理员会帮你补上正式的单位记录。</p>
        <p><strong>研究关键词</strong>这一项特别重要，别的用户搜专家的时候就是靠这个匹配到你，尽量填准确、全面。新增好的档案一开始是"草稿"状态，只有你自己能看到、能改。</p>

        <h4>第二步：补充成果和履历</h4>
        <p>点进某位专家的详情页，里面有几个分区，分别用来补充：</p>
        <ul>
          <li><strong>研究成果</strong>——发表的论文、出版的著作、参与的课题等。成果的级别（国家级/省部级等）会影响这位专家在排序里靠不靠前，如实填写。</li>
          <li><strong>决策影响记录</strong>——成果被领导批示、被政府采纳、写进政策文件这类情况，是本系统区别于普通专家名录的核心内容，级别越高、加分越多。</li>
          <li><strong>学术兼职</strong>——在学术团体、政府咨询机构里担任的职务。</li>
          <li><strong>标签</strong>——给专家打上学科、关键词、政策领域、擅长区域、研究方法这些标签，直接影响搜索时能不能被准确找到，越贴切越好。</li>
        </ul>

        <h4>批量导入</h4>
        <p>如果要一次性录入很多位专家，"专家主档"页面有"批量导入"按钮：先下载表格模板，按格式填好背景信息和研究成果，再上传。表格里有问题的话系统会整体退回、提示哪里错了，不会导入一半错误数据；单位名字如果匹配不上现有单位也不影响导入，只是会在结果里提醒一下，之后再补关联。</p>

        <h4>第三步：提交审核</h4>
        <p>信息填得差不多之后，在专家列表或者"审核台-我发起的"里点"提交审核"。提交后状态会经过下面这几个阶段：</p>
        <div class="status-flow">
          <el-tag type="info">草稿</el-tag>
          <el-icon><Right /></el-icon>
          <el-tag type="warning">待单位审核</el-tag>
          <el-icon><Right /></el-icon>
          <el-tag type="warning">待市级审核</el-tag>
          <el-icon><Right /></el-icon>
          <el-tag type="success">已发布</el-tag>
        </div>
        <p style="margin-top: 10px;">如果中途被退回（单位或市级都可能退回），状态会变成"已退回"，同时会附上退回原因，看完原因改好之后可以重新提交，走的还是这几步。<strong>只有"已发布"状态的专家才会出现在搜索结果里。</strong></p>
      </el-tab-pane>

      <el-tab-pane v-if="isOrgReviewer || isCityReviewer || isAdmin" label="怎么审核" name="reviewer">
        <p>在"审核台"的"待我审核的"里能看到需要你处理的记录。单位审核员看到的是本单位提交的；市级审核员看到的是已经通过单位审核、等待最终确认的。</p>
        <p>每一条记录可以：</p>
        <ul>
          <li><strong>通过</strong>——单位审核通过之后会转给市级审核；市级审核通过之后就正式收录、可以被搜到了。</li>
          <li><strong>退回</strong>——需要填一句退回原因，对方会看到这句话，方便知道要改哪里。</li>
          <li><strong>查看审核记录</strong>——完整看一遍这条记录经历过的每一步。</li>
        </ul>

        <h4>一次处理很多条</h4>
        <p>如果要审核的记录很多，可以把想通过的记录都勾上，点"批量通过"，确认一下就能一次全部处理完，不用一条条点。这个操作会跳过逐条查看，建议只对确实没问题、不需要单独退回的记录使用。</p>

        <h4>"免审核发布"标签</h4>
        <p>有些记录的状态旁边会带一个橙色的"免审核发布"标签，意思是这条记录是系统管理员统一导入历史数据时特批直接收录的，没有真的走完审核流程——这是正常情况，不用当成错误处理。</p>

        <h4 v-if="isOrgReviewer">给本单位的专家开账号</h4>
        <p v-if="isOrgReviewer">"本单位账号管理"里可以直接给本单位的专家开通账号，不用麻烦系统管理员，开出来的账号权限固定就是"自己维护自己的资料"，不会误开成别的权限。</p>
      </el-tab-pane>

      <el-tab-pane label="怎么找专家" name="search">
        <p>"专家检索推荐"这个页面就是用来找人的：输入你关心的主题（比如"防汛方案"），也可以按学科、擅长区域、职称再筛一下，点检索。</p>
        <p>只有正式收录（已发布）、并且确实跟你输入的主题相关的专家才会出现在结果里——不会因为职称高就硬塞一个完全不相关的人进来。</p>
        <p>结果列表会显示这位专家跟你搜索主题的匹配程度、过往成果情况、决策影响力、社会贡献这些信息，综合排出一个先后顺序，越靠前越值得优先考虑。</p>
        <p>如果想要一份名单带走存档或者给别人看，点"导出当前结果"，会把当前搜到的全部结果导出成一份 Excel 表格。</p>
      </el-tab-pane>

      <el-tab-pane v-if="isOrgReviewer || isCityReviewer || isAdmin" label="数据统计" name="dashboard">
        <p>"统计概览"页面能一眼看到：一共收录了多少位专家、还有多少条记录在等审核、各个审核阶段分别有多少、最近一段时间新增了多少、专家都集中在哪些单位和学科。单位审核员看到的统计只包含本单位的数据。</p>
      </el-tab-pane>

      <el-tab-pane v-if="isAdmin" label="系统设置" name="admin">
        <h4>账号管理</h4>
        <p>在"用户管理"里给账号分配角色（个人申报人/单位审核员/市级审核员）和所属单位，这两项都要设置，不然对方登录后没法正常使用——比如没设单位，申报人提交的东西就没有归属，审核员也看不到自己该审的记录。</p>

        <h4>单位管理</h4>
        <p>"单位管理"页面维护系统里有哪些单位，专家档案里"所在单位"的下拉选项、审核该归到哪个单位、统计报表里的单位分布，都是以这里为准。</p>

        <h4>调整评分标准</h4>
        <p>在"字典管理"里能找到几个专门给专家库用的配置项，改完立刻生效，不用重启系统：</p>
        <ul>
          <li><strong>专家库-成果级别权重</strong> —— 不同级别的成果各自值多少分</li>
          <li><strong>专家库-采纳单位级别权重</strong> —— 决策被不同级别单位采纳，各自值多少分</li>
          <li><strong>专家库-职称权重</strong> —— 不同职称各自值多少分</li>
          <li><strong>专家库-社会贡献权重</strong> —— 学术兼职、荣誉称号各自值多少分</li>
          <li><strong>专家库-综合排序权重</strong> —— 最终排序时，成果、决策影响、职称、社会贡献这四项各占多大比重</li>
        </ul>

        <h4>审核开关</h4>
        <p>同样在"字典管理"里，有一项叫<strong>专家库-审核开关</strong>。正常情况下这个开关是打开的，新增/导入/提交的专家都要走完整的审核流程。如果要一次性批量导入一大批已经确认过的历史数据，可以临时把这个开关关掉，这样新增和导入的记录会跳过审核直接收录，处理完记得改回打开状态。开关关闭期间收录的记录会带上"免审核发布"标签，方便以后查证。</p>

        <h4>一次处理大批数据</h4>
        <p>审核台的"批量通过"管理员也能用，两个审核阶段的队列都能看、都能批量处理，适合集中处理大批历史数据的场景。</p>
      </el-tab-pane>

      <el-tab-pane label="常见问题" name="faq">
        <div class="faq-item"><div class="q">提交审核之后想改怎么办？</div><div class="a">正在审核中的记录不能直接改。可以先联系审核员把它退回来，退回之后就能编辑了，改完重新提交。</div></div>
        <div class="faq-item"><div class="q">为什么搜不到某位专家？</div><div class="a">先确认这位专家的档案是不是"已发布"状态，没发布就搜不到；再检查搜索的关键词是不是跟这位专家的研究方向、成果、标签相关。</div></div>
        <div class="faq-item"><div class="q">忘记密码怎么办？</div><div class="a">找系统管理员帮忙重置，重置后用默认密码登录，登录后自己再改成新密码。</div></div>
        <div class="faq-item"><div class="q">单位审核员能看到别的单位的专家吗？</div><div class="a">看不到。系统是按单位自动分开的，每个单位审核员只能看到、只能审核自己单位提交的记录，这是系统统一控制的，没有办法绕过去。</div></div>
        <div class="faq-item"><div class="q">所在单位下拉列表里找不到我的单位怎么办？</div><div class="a">可以先直接输入单位名字保存，不影响正常使用，之后找系统管理员把这个单位正式补进单位列表、重新关联一下就行。</div></div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({
  name: 'ExpertHelp'
})

const userStore = useUserStore()
const authorityId = computed(() => userStore.userInfo.authorityId)
const isAdmin = computed(() => authorityId.value === 1)
const isApplicant = computed(() => authorityId.value === 9001)
const isOrgReviewer = computed(() => authorityId.value === 9002)
const isCityReviewer = computed(() => authorityId.value === 9003)

const myRoleKey = computed(() => {
  if (isAdmin.value) return 'admin'
  if (isApplicant.value) return 'applicant'
  if (isOrgReviewer.value) return 'org'
  if (isCityReviewer.value) return 'city'
  return ''
})

const roleCards = [
  { key: 'applicant', name: '个人申报人', icon: 'User', desc: '填报自己的专家资料、成果和履历，提交审核' },
  { key: 'org', name: '单位审核员', icon: 'CircleCheck', desc: '审核本单位提交的专家资料，给本单位开申报账号' },
  { key: 'city', name: '市级审核员', icon: 'Select', desc: '对已通过单位审核的资料做最终审核，通过后正式收录' },
  { key: 'admin', name: '管理员', icon: 'Setting', desc: '管理账号、单位、评分标准，处理全部审核事务' }
]

const activeTab = ref(
  isApplicant.value ? 'applicant' : (isOrgReviewer.value || isCityReviewer.value) ? 'reviewer' : isAdmin.value ? 'admin' : 'login'
)
</script>

<style scoped>
.expert-help {
  padding: 4px;
}
.hero {
  background: linear-gradient(135deg, var(--el-color-primary-light-9), var(--el-fill-color-light));
  border-radius: 12px;
  padding: 24px 28px;
  margin-bottom: 20px;
}
.hero h2 {
  margin: 0 0 8px;
}
.hero p {
  margin: 0;
  color: var(--el-text-color-regular);
  line-height: 1.8;
}
.role-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 14px;
  margin-bottom: 20px;
}
.role-card {
  border: 1px solid var(--el-border-color);
  border-radius: 10px;
  padding: 16px;
  transition: all 0.2s;
}
.role-card.active {
  border-color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
}
.role-icon {
  font-size: 26px;
  color: var(--el-color-primary);
  margin-bottom: 8px;
}
.role-name {
  font-weight: 600;
  margin-bottom: 6px;
}
.role-desc {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  line-height: 1.6;
}
.help-tabs {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 10px;
  padding: 8px 20px 20px;
}
h4 {
  margin: 16px 0 8px;
}
p {
  line-height: 1.8;
  margin: 8px 0;
}
ul {
  line-height: 1.9;
  padding-left: 20px;
}
.status-flow {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.status-flow .el-icon {
  color: var(--el-text-color-placeholder);
}
.faq-item {
  margin-bottom: 16px;
}
.faq-item .q {
  font-weight: 600;
  margin-bottom: 4px;
}
.faq-item .q::before {
  content: "Q  ";
  color: var(--el-color-primary);
}
.faq-item .a {
  color: var(--el-text-color-regular);
  line-height: 1.8;
}
.faq-item .a::before {
  content: "A  ";
  color: var(--el-color-success);
  font-weight: 600;
}
</style>
