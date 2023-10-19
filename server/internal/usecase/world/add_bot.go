package usecase_world

func (uc *Usecase) AddBot() error {
	b, err := uc.world.NewBot()
	if err != nil {
		return err
	}

	uc.ncProducer.OnPlayerConnect(b.GetPlayer().ID())
	uc.ncProducer.SpawnPlayer(b.GetPlayer())
	return nil
}
