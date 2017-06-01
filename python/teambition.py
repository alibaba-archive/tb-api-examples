#! -*- coding: utf-8 -*-
"""
   a simple teambition API SDK for python
"""
import six
import requests

if six.PY2:
    from urllib import urlencode
    from urlparse import urljoin
else:
    from urllib.parse import urlencode
    from urllib.parse import urljoin

HOST = 'api.teambition.com'
AUTH_HOST = 'account.teambition.com'


class Teambition(requests.Session):
    """
    a Teambition SDK

    Basic Usage::

      >>> from teambition import Teambition
      >>> sdk = Teambition(client_id, client_secret)
      >>> sdk.get_authorize_url(redirect_url)  # get teambition authorize url

    Access to authorize url with brower Get the redirect `code` . then::

      >>> access_token = sdk.get_access_token(code)
      >>> sdk.set_token(token)
      >>> sdk.get('users/me')
    """

    def __init__(self, client_id, client_secret, protocol='https'):
        """
        :class:`Teambtion`.
        
        :param client_id: :class:`string` teambition client id
        :param client_secret: :class:`string` teambition client secret
        """
        self._client_id = client_id
        self._client_secret = client_secret

        self._root = protocol + '://' + HOST
        self._auth_root = protocol + '://' + AUTH_HOST
        
        super(Teambition, self).__init__()

        self.headers.update({
            'Content-Type': 'application/json'
        })

    def set_token(self, token):
        """
        set token to session headers.

        :class:`Teambtion`.
        
        :param token: :class:`string` the oauth2 access token.
        """
        self.headers.update({
            'Authorization': 'OAuth2 ' + token
        })

    def get_authorize_url(self, redirect_url):
        """
        give authorize url

        :class:`Teambtion`.
        
        :param redirect_url: :class:`string` the callback url
        :rtype: string
        """
        _url = urljoin(self._auth_root, '/oauth2/authorize')
        params = {
            'client_id': self._client_id,
            'redirect_uri': redirect_url
        }
        return _url + '?' + urlencode(params)

    def get_access_token(self, code, state=None):
        """
        give access token

        :class:`Teambtion`.
        
        :param code: :class:`string` callback code
        :param state: :class:`string` https://auth0.com/docs/protocols/oauth2/oauth-state
        :rtype: string
        """
        payload = {
            'client_id': self._client_id,
            'client_secret': self._client_secret,
            'code': code
        }

        if state:
            payload['state'] = state

        resp = super(Teambition, self).request(
            'post', 
            urljoin(self._auth_root, '/oauth2/access_token'), 
            json=payload
        )
        return resp.json()["access_token"]

    def request(self, method, path, *args, **kwargs):
        """
        request method
        """
        url = urljoin(self._root, '/api' + path)
        resp = super(Teambition, self).request(
            method, url, *args, **kwargs
        )
        
        resp.raise_for_status()
        return resp.json()
