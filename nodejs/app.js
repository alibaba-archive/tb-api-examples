'use strict'

const Koa = require('koa')
const router = require('koa-router')()
const Teambition = require('teambition')
const app = module.exports = new Koa()

let sdk = new Teambition()

const CLIENT_ID = 'c9c44aa0-45f8-11e7-85e5-25300cc3a657'
const CLIENT_SECRET = 'e297f011-ea56-4421-8be9-6477933e1591'
const REDIRECT_URI = 'http://localhost:3000/tb/callback'

// 跳转至授权地址
router.get('/auth', function * () {
  let gotoAuth = sdk.getAuthorizeUrl(CLIENT_ID, REDIRECT_URI)
  this.redirect(gotoAuth)
})

// 处理授权回调
router.get('/tb/callback', sdk.authCoCallback(CLIENT_ID, CLIENT_SECRET), function * () {
  let tbCallbackBody = this.request.callbackBody
  if (tbCallbackBody && tbCallbackBody.access_token) {
    sdk.token = tbCallbackBody.access_token

    // 创建一个项目
    let project = yield sdk.post('/projects', { name: '示例项目 by api' })

    // 获取项目的任务列表
    let tasklists = yield sdk.get(`/projects/${project._id}/tasklists`)

    // 创建任务
    let task = yield sdk.post('/tasks', {
      _projectId: project._id,
      _tasklistId: tasklists[0]._id,
      _stageId: tasklists[0].hasStages[0]._id,
      content: '示例任务 by api'
    })

    // 完成任务
    yield sdk.put(`/tasks/${task._id}/isDone`, { isDone: true })

    // 创建日程
    let event = yield sdk.post('/events', {
      title: '示例日程标题 by api',
      content: '示例日程备注 by api',
      startDate: new Date(2017, 5, 1),
      endDate: new Date(2017, 5, 4),
      _projectId: project._id
    })

    // 更新日程的标题
    yield sdk.put(`/events/${event._id}`, { title: '新示例日程标题 by api' })

    // 邀请新成员到项目
    let members = yield sdk.post(`/projects/${project._id}/members`, { email: 'teambition@teambition.com' })

    // 移除成员
    yield sdk.del(`/members/${members[0]._memberId}`)

    // 删除日程
    yield sdk.del(`/events/${event._id}`)

    // 删除任务
    yield sdk.del(`/tasks/${task._id}`)

    // 删除项目
    yield sdk.del(`/projects/${project._id}`)

    this.body = {
      project: project,
      tasklists: tasklists,
      task: task,
      event: event,
      members: members
    }
  } else {
    this.body = tbCallbackBody
  }
})

app.use(router.routes())

app.listen(3000)
