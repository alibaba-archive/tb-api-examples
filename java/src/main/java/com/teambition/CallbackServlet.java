package com.teambition;

import okhttp3.*;
import okhttp3.Request.Builder;
import org.json.JSONArray;
import org.json.JSONObject;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

/**
 * @author Xu Jingxin
 */
public class CallbackServlet extends HttpServlet {
    private static OkHttpClient client = new OkHttpClient();
    private static MediaType mediaType = MediaType.parse("application/json");

    @Override
    protected void doGet (HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {
        String code = request.getParameter("code");
        response.setContentType("application/json; charset=utf-8");
        if (code == null) {
            response.setStatus(HttpServletResponse.SC_FORBIDDEN);
            response.getWriter().println(new JSONObject().put("error", "code 参数错误")
                    .toString());
            return;
        }
        // 获取 access_token
        String accessToken = getAccessToken(code);
        // 创建项目
        JSONObject project = createProject(accessToken);
        // 获取项目中的任务列表
        JSONArray tasklists = getTasklists(accessToken, project.getString("_id"));
        // 创建任务
        JSONObject task = createTask(accessToken, tasklists.getJSONObject(0), project);
        // 完成任务
        accomplishTask(accessToken, task.getString("_id"));
        // 创建日程
        JSONObject event = createEvent(accessToken, project.getString("_id"));
        // 更新日程
        updateEvent(accessToken, event.getString("_id"));
        // 邀请成员
        JSONArray members = inviteMember(accessToken, project.getString("_id"));
        // 移除成员
        removeMember(accessToken, members.getJSONObject(0).getString("_id"));
        // 移除日程
        removeEvent(accessToken, event.getString("_id"));
        // 移除任务
        removeTask(accessToken, task.getString("_id"));
        // 移除项目
        removeProject(accessToken, project.getString("_id"));

        JSONObject data = new JSONObject().put("project", project)
                .put("tasklists", tasklists)
                .put("tasks", task)
                .put("event", event)
                .put("members", members);
        response.setStatus(HttpServletResponse.SC_OK);
        response.getWriter().println(data.toString());
    }

    /**
     * 获取 access_token
     *
     * @param code
     * @return
     * @throws IOException
     */
    private String getAccessToken (String code) throws IOException {
        JSONObject payload = new JSONObject()
                .put("client_id", ApiConst.CLIENT_ID)
                .put("client_secret", ApiConst.CLIENT_SECRET)
                .put("code", code);
        RequestBody body = RequestBody.create(mediaType, payload.toString());
        JSONObject data = (JSONObject) request(new Request.Builder()
                .url(ApiConst.AUTH_HOST + "/oauth2/access_token")
                .post(body));
        return data.get("access_token").toString();
    }

    /**
     * 不带 token 请求
     *
     * @param builder
     * @return
     * @throws IOException
     */
    private Object request (Builder builder) throws IOException {
        Request request = builder
                .addHeader("content-type", "application/json")
                .build();
        Response response = client.newCall(request).execute();
        String resData = response.body().string();
        if (resData.charAt(0) == '[') {
            return new JSONArray(resData);
        } else {
            return new JSONObject(resData);
        }
    }

    /**
     * 带 token 请求
     *
     * @param builder
     * @param accessToken
     * @return
     * @throws IOException
     */
    private Object request (Builder builder, String accessToken) throws IOException {
        Request request = builder
                .addHeader("content-type", "application/json")
                .addHeader("authorization", "OAuth2 " + accessToken)
                .build();
        Response response = client.newCall(request).execute();
        String resData = response.body().string();
        if (resData.charAt(0) == '[') {
            return new JSONArray(resData);
        } else {
            return new JSONObject(resData);
        }
    }

    /**
     * 创建项目
     *
     * @param accessToken
     * @return
     * @throws IOException
     */
    private JSONObject createProject (String accessToken) throws IOException {
        JSONObject payload = new JSONObject()
                .put("name", "示例项目");
        RequestBody body = RequestBody.create(mediaType, payload.toString());
        return (JSONObject) request(new Request.Builder()
                .url(ApiConst.API_HOST + "/projects")
                .post(body), accessToken);
    }

    /**
     * 获取任务列表
     *
     * @param accessToken
     * @return
     * @throws IOException
     */
    private JSONArray getTasklists (String accessToken, String _projectId) throws IOException {
        return (JSONArray) request(new Request.Builder()
                .url(ApiConst.API_HOST + "/projects/" + _projectId + "/tasklists")
                .get(), accessToken);
    }

    /**
     * 创建任务
     */
    private JSONObject createTask (String accessToken, JSONObject tasklist, JSONObject project) throws IOException {
        JSONObject payload = new JSONObject()
                .put("content", "示例任务 by api")
                .put("_projectId", project.getString("_id"))
                .put("_tasklistId", tasklist.getString("_id"));
        JSONArray stages = tasklist.getJSONArray("hasStages");
        payload.put("_stageId", stages.getJSONObject(0).getString("_id"));
        RequestBody body = RequestBody.create(mediaType, payload.toString());
        return (JSONObject) request(new Request.Builder()
                .url(ApiConst.API_HOST + "/tasks")
                .post(body), accessToken);
    }

    /**
     * 完成任务
     */
    private JSONObject accomplishTask (String accessToken, String _taskId) throws IOException {
        JSONObject payload = new JSONObject()
                .put("isDone", true);
        RequestBody body = RequestBody.create(mediaType, payload.toString());
        return (JSONObject) request(new Request.Builder()
                .url(ApiConst.API_HOST + "/tasks/" + _taskId + "/isDone")
                .put(body), accessToken);
    }

    /**
     * 创建日程
     */
    private JSONObject createEvent (String accessToken, String _projectId) throws IOException {
        JSONObject payload = new JSONObject()
                .put("_projectId", _projectId)
                .put("title", "示例日程标题 by api")
                .put("content", "示例日程备注 by api")
                .put("startDate", "2017-06-01T14:50:15.035Z")
                .put("endDate", "2017-06-05T18:00:00.035Z");
        RequestBody body = RequestBody.create(mediaType, payload.toString());
        return (JSONObject) request(new Request.Builder()
                .url(ApiConst.API_HOST + "/events")
                .post(body), accessToken);
    }

    /**
     * 更新日程标题
     */
    private JSONObject updateEvent (String accessToken, String _eventId) throws IOException {
        JSONObject payload = new JSONObject()
                .put("title", "新示例日程标题 by api");
        RequestBody body = RequestBody.create(mediaType, payload.toString());
        return (JSONObject) request(new Request.Builder()
                .url(ApiConst.API_HOST + "/events/" + _eventId)
                .put(body), accessToken);
    }

    /**
     * 邀请成员
     */
    private JSONArray inviteMember (String accessToken, String _projectId) throws IOException {
        JSONObject payload = new JSONObject()
                .put("email", "test@teambition.com");
        RequestBody body = RequestBody.create(mediaType, payload.toString());
        return (JSONArray) request(new Request.Builder()
                .url(ApiConst.API_HOST + "/projects/" + _projectId + "/members")
                .post(body), accessToken);
    }

    /**
     * 移除成员
     */
    private JSONObject removeMember (String accessToken, String _memberId) throws IOException {
        return (JSONObject) request(new Request.Builder()
                .url(ApiConst.API_HOST + "/members/" + _memberId)
                .delete(), accessToken);
    }

    /**
     * 删除日程
     */
    private JSONObject removeEvent (String accessToken, String _eventId) throws IOException {
        return (JSONObject) request(new Request.Builder()
                .url(ApiConst.API_HOST + "/events/" + _eventId)
                .delete(), accessToken);
    }

    /**
     * 删除任务
     */
    private JSONObject removeTask (String accessToken, String _taskId) throws IOException {
        return (JSONObject) request(new Request.Builder()
                .url(ApiConst.API_HOST + "/tasks/" + _taskId)
                .delete(), accessToken);
    }

    /**
     * 删除项目
     */
    private JSONObject removeProject (String accessToken, String _projectId) throws IOException {
        return (JSONObject) request(new Request.Builder()
                .url(ApiConst.API_HOST + "/projects/" + _projectId)
                .delete(), accessToken);
    }
}
