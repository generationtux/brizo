import { Injectable }       from '@angular/core';
import { tokenNotExpired, JwtHelper }  from 'angular2-jwt';
import { Headers, RequestOptions } from '@angular/http';

@Injectable()
export class AuthService {
  private jwt: JwtHelper;

  constructor() {
    this.jwt = new JwtHelper();
  }

  validateJwtToken() {
    return tokenNotExpired();
  }

  getToken() :string {
    return localStorage.getItem('id_token');
  }

  clearToken() :void {
    localStorage.removeItem('id_token');
  }

  jwtRequestOptions() :RequestOptions {
    let headers = new Headers({'Authorization': 'Bearer ' + this.getToken()})
    return new RequestOptions({headers})
  }

  user() :User | null {
    const token = this.getToken();
    if (token === null) {
      return null;
    }

    return new User(this.jwt.decodeToken(token));
  }
}

export class User {
  public name: string;
  public email: string;

  constructor(props: any = {}) {
    this.name = props.name || "";
    this.email = props.email || "";
  }

  public image() :string {
    let md5 = require('js-md5');
    let hash = md5(this.email);
    return "https://www.gravatar.com/avatar/" + hash + "?s=25";
  }
}
