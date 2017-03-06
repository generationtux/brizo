import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from '@angular/router';
import { SharedModule } from '../shared.module';

import { AuthGuard } from '../auth/auth.guard';
import { HelpComponent } from './help.component';
import { HelpAccessTokensComponent } from './components/help-access-tokens.component';
import { HelpAppsComponent } from './components/help-apps.component';
import { HelpEnvironmentsComponent } from './components/help-environments.component';
import { HelpVersionsComponent } from './components/help-versions.component';
import { HelpViewUsersComponent } from './components/help-view-users.component';
import { HelpInviteUsersComponent } from './components/help-invite-users.component';

const helpRoutes: Routes = [
  { path: 'help', component: HelpComponent, canActivate: [AuthGuard] },
  { path: 'help/apps', component: HelpAppsComponent, canActivate: [AuthGuard] },
  { path: 'help/environments', component: HelpEnvironmentsComponent, canActivate: [AuthGuard] },
  { path: 'help/versions', component: HelpVersionsComponent, canActivate: [AuthGuard] },
  { path: 'help/access-tokens', component: HelpAccessTokensComponent, canActivate: [AuthGuard] },
  { path: 'help/view-users', component: HelpViewUsersComponent, canActivate: [AuthGuard] },
  { path: 'help/invite-users', component: HelpInviteUsersComponent, canActivate: [AuthGuard] },
];

@NgModule({
  imports: [
    CommonModule,
    RouterModule.forChild(helpRoutes),
    SharedModule,
  ],
  declarations: [
    HelpComponent,
    HelpAccessTokensComponent,
    HelpAppsComponent,
    HelpEnvironmentsComponent,
    HelpVersionsComponent,
    HelpInviteUsersComponent,
    HelpViewUsersComponent,
  ],
  exports: [
    HelpComponent,
    HelpAccessTokensComponent,
    HelpAppsComponent,
    HelpEnvironmentsComponent,
    HelpVersionsComponent,
    HelpViewUsersComponent,
    HelpInviteUsersComponent,
    RouterModule,
  ],
})

export class HelpModule {}
