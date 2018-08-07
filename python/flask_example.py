#! -*- coding: utf-8 -*-
from flask import Flask, redirect, url_for,\
     abort, request, jsonify
from teambition import Teambition

app = Flask(__name__)

app.config.update(
    SERVER_NAME = 'localhost:3000'
)

CLIENT_ID = 'xxxxxx'
CLIENT_SECRET = 'xxxxxx'


tb = Teambition(CLIENT_ID, CLIENT_SECRET)

@app.route('/tb/auth')
def auth():
    """
    跳转Teambition授权地址
    """
    return redirect(tb.get_authorize_url(
        url_for('callback', _external=True)
    ))


@app.route('/tb/callback')
def callback():
    """
    完成授权地址确认后，Teambition API跳转到此endpoint
    """
    code = request.args.get('code', None)

    if not code:
        abort(500)

    # 通过返回code， 获取AccessToken
    access_token = tb.get_access_token(code)

    # set token
    tb.set_token(access_token)

    # 创建一个项目
    project = tb.post('/projects', json={'name': '示例项目 by api'})

    # 获取项目的任务列表
    tasklists = tb.get('/projects/%s/tasklists' % project["_id"])

    # 创建任务
    task = tb.post('/tasks', json={
        '_projectId': project["_id"],
        '_tasklistId': tasklists[0]['_id'],
        '_stageId': tasklists[0]['hasStages'][0]['_id'],
        'content': '示例任务 by api'
    })

    # 完成任务
    tb.put("/tasks/%s/isDone" % task['_id'], json={'isDone': True})

    # 创建日程
    event = tb.post('/events', json={
        'title': '示例日程标题 by api',
        'content': '示例日程备注 by api',
        'startDate': '2017-06-01T14:50:15.035Z',
        'endDate': '2017-06-05T18:00:00.035Z',
        '_projectId': project['_id']
    })

    # 更新日程的标题
    tb.put('/events/%s' % event['_id'], 
        json={'title': '新示例日程标题 by api'})
    
    # 邀请新成员到项目
    members = tb.post('/projects/%s/members' % project['_id'], 
        json={'email': 'teambition@teambition.com' })

    # 移除成员
    tb.delete('/members/%s' % members[0]['_memberId'])

    # 删除日程
    tb.delete('/events/%s' % event['_id'])

    # 删除任务
    tb.delete('/tasks/%s' % task['_id'])

    # 删除项目
    tb.delete('/projects/%s' % project['_id'])

    return jsonify({
        'project': project,
        'tasklists': tasklists,
        'task': task,
        'event': event,
        'members': members
    })


if __name__ == '__main__':
    app.run()
