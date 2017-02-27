import { Observable } from 'rxjs/Rx';
import { Component, Injectable } from '@angular/core';
import { Http, Response, Request, RequestMethod } from '@angular/http';
import { FormGroup, FormControl, Validators, FormBuilder } from '@angular/forms';

import { AuthService } from '../auth/auth.service';

@Component({
    selector:       'access-tokens-form',
    // template:       '<div>Hello!</div>',
    templateUrl:    './access-tokens.html',
})

@Injectable()
export class AccessTokenComponent {
  private accessTokensIndexUrl = '/api/v1/access-tokens';
  private accessTokensCreateUrl = '/api/v1/access-tokens';
  public form: FormGroup;
  public tokens: Observable<AccessToken[]>;

  constructor(private formBuilder: FormBuilder, private http: Http, private auth: AuthService) {
    this.form = formBuilder.group({});
  }

  ngOnInit() {
    this.hydrateAccessTokens()
  }

  submitAccessTokenForm() {
    this.http.post(this.accessTokensCreateUrl, null, this.auth.jwtRequestOptions()).subscribe(
      (res: Response) => {
        this.hydrateAccessTokens();
      }
    )
  }

  getAccessTokens(): Observable<AccessToken[]> {
    return this.http.get(this.accessTokensIndexUrl, this.auth.jwtRequestOptions())
      .map((res: Response) => res.json())
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }

  private hydrateAccessTokens() {
    this.tokens = this.getAccessTokens()
  }
}

export class AccessToken {
    constructor(
        public token: string,
        public created_at: string,
    ){}
}
