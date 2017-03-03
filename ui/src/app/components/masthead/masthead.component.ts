import { Component } from '@angular/core';
import { Router } from '@angular/router';

import { AuthService, User } from '../../modules/auth/auth.service';

@Component({
    selector:       'masthead',
    templateUrl:    './masthead.html',
})

export class MastheadComponent {
  public user: User | null;

  constructor(private router: Router, private auth: AuthService) {
    this.loadUser();
  }

  loadUser(): void {
    this.user = this.auth.user();
  }

  logout(event: any) {
    event.preventDefault();
    // @todo move to auth service
    localStorage.removeItem('id_token');
    this.router.navigate(['login'])
  }
}
