import { Component } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
    selector:       'login',
    styleUrls:      ['./login.css'],
    templateUrl:    './login.html'
})

export class LoginComponent {
  private loginError :string

  constructor(private router: Router, private route: ActivatedRoute) {
    route.queryParams.subscribe(
      (queryParam: any) => {
        let errorCode = queryParam['err'];
        if (errorCode == "1") {
          this.loginError = "Something went wrong"
        } else if (errorCode == "2") {
          this.loginError = "You are not authorized to login."
        }
      }
    );
  }
}
