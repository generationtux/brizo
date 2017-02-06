import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from '@angular/router';
import { AuthModule } from '../auth/auth.module';
import { LoginComponent } from './login.component';

import { AuthService } from '../auth/auth.service';

const loginRoutes: Routes = [
  { path: 'login', component: LoginComponent },
];

@NgModule({
  imports:      [ AuthModule, CommonModule, RouterModule.forChild(loginRoutes) ],
  declarations: [ LoginComponent ],
  exports:      [ LoginComponent, RouterModule ],
  providers:    [ AuthService ],
})
export class LoginModule { };
