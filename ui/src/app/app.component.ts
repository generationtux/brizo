import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { AuthService } from './modules/auth/auth.service';

@Component({
    selector:       'app',
    templateUrl:    './app.html'
})

export class AppComponent {

  constructor(private auth: AuthService) {}

  isLoggedIn(): boolean {
    return this.auth.user() !== null;
  }
}
