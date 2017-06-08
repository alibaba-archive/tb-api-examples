package com.teambition;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.net.URLEncoder;

/**
 * @author Xu Jingxin
 */
public class AuthServlet extends HttpServlet {

    @Override
    protected void doGet (HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {
        response.sendRedirect(getAuthUrl());
    }

    /**
     * 获取临时 code
     *
     * @return String
     */
    private String getAuthUrl () throws UnsupportedEncodingException {
        String authUrl = ApiConst.AUTH_HOST + "/oauth2/authorize";
        authUrl += "?client_id=" + ApiConst.CLIENT_ID +
                "&client_secret=" + ApiConst.CLIENT_SECRET +
                "&redirect_uri=" + URLEncoder.encode("http://localhost:3000/tb/callback", "UTF-8");
        return authUrl;
    }
}
