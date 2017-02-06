import { Observable } from 'rxjs/Rx';
import { Component, Injectable } from '@angular/core';
import { Http, Response, Request, RequestMethod } from '@angular/http';
import { FormGroup, FormControl, Validators, FormBuilder } from '@angular/forms';

import { AuthService } from '../auth/auth.service';

@Component({
    selector:       'invite-form',
    templateUrl:    './invite.html'
})

@Injectable()
export class InviteComponent {
  private invitesGetUrl = '/api/v1/users/invites'
  private invitePostUrl = '/api/v1/users';
  public form: FormGroup;
  public invitees: Observable<Invitee[]>;
  public username = new FormControl(null, Validators.required);

  constructor(private formBuilder: FormBuilder, private http: Http, private auth: AuthService) {
    this.form = formBuilder.group({
      'username': this.username,
    });
  }

  ngOnInit() {
    this.hydrateInvitees()
  }

  submitInviteForm() {
    const data = {
      'username': this.username.value,
    }
    this.http.post(this.invitePostUrl, data, this.auth.jwtRequestOptions()).subscribe(
      (res: Response) => {
        this.hydrateInvitees();
        this.username.setValue(null);
      }
    )
  }

  getInvitees(): Observable<Invitee[]> {
    return this.http.get(this.invitesGetUrl, this.auth.jwtRequestOptions())
      .map((res: Response) => res.json())
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }

  private hydrateInvitees() {
    this.invitees = this.getInvitees()
  }

  private validateInvitee(username: FormControl) {
    // @todo
  }
}

export class Invitee {
    constructor(
        public username: string,
    ){}
}
