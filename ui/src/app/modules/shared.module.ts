import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { MastheadComponent } from '../components/masthead/masthead.component';

import { ApplicationService } from './applications/application.service';
import { AuthService } from './auth/auth.service';

@NgModule({
  imports:      [ CommonModule, FormsModule, ReactiveFormsModule, RouterModule ],
  declarations: [ MastheadComponent ],
  exports:      [ MastheadComponent ],
  providers:    [ ApplicationService, AuthService ],
})

export class SharedModule {}
